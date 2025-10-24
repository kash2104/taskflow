package grpcClient

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	pb "github.com/kash2104/taskflow/master/pb"
	taskhandler "github.com/kash2104/taskflow/worker/task-handler"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)


const MAXTHREADS = 10;
func StartGrpcClient(address string){

	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()));

	if err != nil{
		log.Fatalf("Could not connect to grpc server %v", err);
	}
	defer conn.Close()

	client := pb.NewExecuteTaskClient(conn);
	semaphore := make(chan struct{}, MAXTHREADS);

	var wg sync.WaitGroup;

	for{
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second);

		execTask,err := client.GetExecutionRequest(ctx,&pb.ExecutionRequest{})


		cancel();

		if err != nil{
			log.Printf("No task to fetch in the queue");
			time.Sleep(5*time.Second);
			log.Printf("Waiting for tasks to appear in the queue")
			continue;
		}

		fmt.Println("Acquiring semaphore");
		semaphore <- struct{}{}

		wg.Add(1);

		go func(task *pb.ExecutionResponse){
			defer func(){
				<- semaphore;
				wg.Done();
			}()

			//need to make a function that will execute the task
			taskhandler.ProcessTaskSimultaneously(task,client);
			
			
			// fmt.Println("Task received is: ", task);
			// id,_ := primitive.ObjectIDFromHex(task.Id);
			// stdout,stderr,err := sandbox.RunInSandbox(id,task.Language,task.Code);
			
			// fmt.Println("output is:",stdout);
			// fmt.Println("errout is: ", stderr);
			// fmt.Println("actual error: ", err);
		}(execTask)



	}
}