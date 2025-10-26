package taskhandler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/kash2104/taskflow/master/pb"
	"github.com/kash2104/taskflow/worker/sandbox"
)

var timelimit int = 60
var memoryLimit int = 256 * 1024 * 1024
var Client pb.ExecuteTaskClient;

func RunAndCompute(ctx context.Context, task *pb.ExecutionResponse)(bool, error){

	newTask := NewTask(
		task.Id,
		task.Code,
		task.Language,
		timelimit,
		memoryLimit,
	)

	errorChannel := make(chan error, 1);
	completed := make(chan bool, 1);
	taskId := ConvertStringToObjectId(task.Id);

	go func(){
		defer close(completed);

		result,err := RunTask(newTask,ctx); 
		if err != nil{
			errorChannel <- errors.Join(errors.New("error running the code"), err);
			return;
		}

		_, err = Client.UpdateTaskStatus(ctx, &pb.UpdateRequest{
				Id: task.Id,
				Status: "Completed",
		})
		if err != nil{
			log.Fatalf("Failed to update the database after completion %v", err);
			errorChannel <- errors.New("Failed to updated completed status in db");
			return;
		}

		_, err = Client.UpdateTaskResult(ctx, &pb.UpdateRequest{
			Id: task.Id,
			Status: result,
		})
		if err != nil{
			log.Fatalf("Failed to update the database with result %v", err);
			errorChannel <- errors.New("FAiled to updated result in db");
			return;
		}

		completed <- true;

	}()

	select{
		case <- ctx.Done():
			select{
				case <- completed:
					return true,nil
				default:
					defer sandbox.CleanupSandbox(taskId);
					return false, errors.New("context timed out");
			}

		case err := <- errorChannel:
			if err != nil{
				fmt.Println("Errror", err);
				return false, err;
			}
	}
	return false, nil
}

func ProcessTaskSimultaneously(task *pb.ExecutionResponse, client pb.ExecuteTaskClient){

	taskId := ConvertStringToObjectId(task.Id);

	Client = client;
	//set context with timeout
	ctx, cancel := context.WithTimeout(context.Background(),time.Duration(timelimit)*time.Second);


	continueStatus, runErr := RunAndCompute(ctx,task);

	if continueStatus{
		cancel();
	}

	if runErr != nil{
		log.Printf("Error executing the task %v", runErr);
		sandbox.CleanupSandbox(taskId);
	}else{
		log.Printf("Done processing task %s", task.Id);
	}

	sandbox.CleanupSandbox(taskId);
	cancel();
}