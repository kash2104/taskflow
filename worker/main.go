package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kash2104/taskflow/master/pb"
	"github.com/kash2104/taskflow/worker/grpcClient"
)

var Client pb.ExecuteTaskClient;

func main(){

	err := godotenv.Load(".env");
	if err != nil{
		// log.Fatalf("failed to load env")
		log.Println("Cannot find env in local");
	}

	ADDRESS := os.Getenv("ADDRESS");
	grpcClient.StartGrpcClient(ADDRESS)
}
