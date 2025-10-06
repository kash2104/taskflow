package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/kash2104/taskflow/master/db"
	"github.com/kash2104/taskflow/master/pb"
	"github.com/kash2104/taskflow/master/queue"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
)


type server struct{
	pb.UnimplementedExecuteTaskServer
}


var taskQueue queue.Queue;

func StartGRPCServer(taskChannel chan db.UserRequest){

	port := os.Getenv("GRPC_PORT")
	lis, err := net.Listen("tcp",port);
	
	if err != nil{
		log.Fatalf("Failed to start the grpc server %v", err);
	}
	
	grpcServer := grpc.NewServer();
	pb.RegisterExecuteTaskServer(grpcServer, &server{});
	
	TaskChannel = taskChannel
	
	// adding the userRequest from channel to queue
	go func(){
		for task := range TaskChannel{
			
			taskQueue.Push(task);
			log.Println("User Request added from channel to queue");
		}
	}()
	
	go func(){
		fmt.Println("grpc server is running on port" + port);
	
		if err := grpcServer.Serve(lis); err != nil{
			log.Fatalf("Failed to serve the grpc server %v", err);
		}
	}()
}

func (s *server) GetExecutionRequest(context context.Context ,req *pb.ExecutionRequest) (*pb.ExecutionResponse, error){

	var executionTask pb.ExecutionResponse
	dummyTask := pb.ExecutionResponse{
		Id: "60d5ec49f1d2c12a4c8b4567",
		Code: "#include <iostream>\nint main() {\n    std::cout << \"Hello, World!\" << std::endl;\n    return 0;\n}",
		Language: "cpp",
		Status: "Completed",
		Result: "No errors",
	}


	if taskQueue.IsEmpty() {
		return &dummyTask, nil
	}else{
		task,_ := taskQueue.Top()
		
		executionTask.Id = task.Id.Hex()
		executionTask.Language = task.Language
		executionTask.Code = task.Code
		executionTask.Status = task.Status
		executionTask.Result = task.Result

		taskQueue.Pop()

		return &executionTask, nil
	}
}
	

func (s * server) UpdateTaskStatus(context context.Context, req *pb.UpdateRequest) (*pb.UpdateResponse, error){

	id := req.Id

	objectId, err := primitive.ObjectIDFromHex(id);

	if err != nil{
		fmt.Println("Error while converting id from string to objectid");
		return nil,err
	}

	result, err := db.UpdateStatus(objectId);
	if err != nil{
		fmt.Println("Error in rpc from the UpdateStatus db call");
		return nil, err;
	}

	fmt.Println("Updated status via rpc call: ", result);

	return nil,nil
}