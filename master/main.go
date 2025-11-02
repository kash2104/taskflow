package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/kash2104/taskflow/master/db"
	"github.com/kash2104/taskflow/master/server"
)


func main(){


	fmt.Println("Starting the program")


	err := godotenv.Load(".env");
	if err != nil {
		// log.Fatalf("error while loading env")
		log.Println("Cannot find env in local")
	}
	MONGO_URI := os.Getenv("MONGODB_URI");

	if err := db.Connect(MONGO_URI); err != nil{
		log.Fatalf("error connecting db %v", err);
	}

	taskChannel := make(chan db.UserRequest, 10)	

	var wg sync.WaitGroup;
	wg.Add(2);

	go func(){
		fmt.Println("Starting 1st thread");
		defer wg.Done();
		server.StartServer(taskChannel);
	}()

	go func(){
		fmt.Println("Starting 2nd thread");
		defer wg.Done()
		server.StartGRPCServer(taskChannel);
	}()

	wg.Wait();
	close(taskChannel);

}