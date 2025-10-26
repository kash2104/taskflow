package taskhandler

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/kash2104/taskflow/master/pb"
	"github.com/kash2104/taskflow/worker/sandbox"
	"github.com/kash2104/taskflow/worker/sandbox/runners"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


type Task struct{
	Id string
	Code string
	Language string
	TimeLimit int
	MemoryLimit int
}

func NewTask(id string, code string, language string,timelimit int, memorylimit int) *Task{

	return &Task{
		Id: id,
		Code: code,
		Language: language,
		TimeLimit: timelimit,
		MemoryLimit: memorylimit,
	}
}

func ConvertStringToObjectId(id string) primitive.ObjectID{
	objectId, err := primitive.ObjectIDFromHex(id);
	if err != nil{
		log.Fatalf("Failed to convert string to objectId %v", err);
	}
	return objectId;
}


func RunTask(task *Task,ctx context.Context) (string, error){

	taskId := ConvertStringToObjectId(task.Id);
	//initialize sandbox
	if err := sandbox.InitSandbox(taskId); err != nil{
		log.Fatalf("Failed to initialise sandbox %v", err);
	}	

	//copy code to the isolate box
	_, err := runners.WriteCodeToSandbox(taskId,task.Code,task.Language);

	if err != nil{
		return "",err;
	}
	defer sandbox.CleanupSandbox(taskId);
	
	switch task.Language{
		case "cpp":

			//compile cpp code
			if errr := runners.CompileCpp(taskId, task.Code); errr != nil{
				_, err := Client.UpdateTaskStatus(ctx,&pb.UpdateRequest{
					Id: task.Id,
					Status: "Compilation error",
				})

				if err != nil{
					log.Fatalf("Failed to update the database %v", err);
				}

				sandbox.CleanupSandbox(taskId);
				return "",errors.New("compilation went wrong");
			}

			//run cpp code
			stdout,stderr, errr := runners.RunCpp(taskId); 
			if errr != nil{
				_, err := Client.UpdateTaskStatus(ctx, &pb.UpdateRequest{
					Id : task.Id,
					Status: "Time limit exceeded",
				})

				if err != nil{
					log.Fatalf("Failed to update the database after tle %v", err);
				}

				return "",errors.New("time task execution failed");
			}
			// _, err := Client.UpdateTaskStatus(ctx, &pb.UpdateRequest{
			// 	Id: task.Id,
			// 	Status: "Completed",
			// })
			// if err != nil{
			// 	log.Fatalf("Failed to update the database after completion %v", err);
			// }
			fmt.Println("Task execution successful");
			// fmt.Println("stdout: ", stdout);
			// fmt.Println("stderr: ", stderr);

			combinedOutput := "STDOUT:\n" + stdout + "\n\nSTDERR:\n" + stderr
			// fmt.Println(combinedOutput);
			return combinedOutput,nil
			

		case "java":
			if errr := runners.CompileJava(taskId, task.Code); errr != nil{
				_, err := Client.UpdateTaskStatus(ctx,&pb.UpdateRequest{
					Id: task.Id,
					Status: "Compilation error",
				})

				if err != nil{
					log.Fatalf("Failed to update the database %v", err);
				}

				return "",errors.New("compilation went wrong");
			}

			stdout,stderr,errr := runners.RunJava(taskId,task.Code);

			if errr != nil{
				_, err := Client.UpdateTaskStatus(ctx,&pb.UpdateRequest{
					Id: task.Id,
					Status: "Time limit exceeded",
				})

				if err != nil{
					log.Fatalf("Failed to update the database after tle %v",err);
				}

				return "", errors.New("time task execution failed");
			}

			fmt.Println("Task execution successful");
			combinedOutput := "STDOUT:\n" + stdout + "\n\nSTDERR:\n" + stderr
			// fmt.Println(combinedOutput);
			return combinedOutput,nil
		case "go":
			stdout,stderr, errr := runners.RunGo(taskId, task.Code);

			if errr != nil{
				_, err := Client.UpdateTaskStatus(ctx, &pb.UpdateRequest{
					Id: task.Id,
					Status: "Time limited exceeded",
				})

				if err != nil{
					log.Fatalf("failed to update the database after tle %v", err);
				}
				return "", errors.New("time task execution failed");
			}

			fmt.Println("task execution successful");
			combinedOutput := "STDOUT:\n" + stdout + "\n\nSTDERR:\n" + stderr

			return combinedOutput, nil

		case "py":
			stdout,stderr, errr := runners.RunPython(taskId, task.Code);

			fmt.Println(errr);
			if errr != nil{
				_, err := Client.UpdateTaskStatus(ctx, &pb.UpdateRequest{
					Id: task.Id,
					Status: "Time limited exceeded",
				})

				if err != nil{
					log.Fatalf("failed to update the database after tle %v", err);
				}
				return "", errors.New("time task execution failed");
			}

			fmt.Println("task execution successful");
			combinedOutput := "STDOUT:\n" + stdout + "\n\nSTDERR:\n" + stderr

			return combinedOutput, nil
	}

	return "Language not identified",nil;
}