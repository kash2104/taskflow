package server

import (
	"fmt"
	"log"
	"time"

	"github.com/kash2104/taskflow/master/db"
)

/*--------------------------Pending Queue-----------------------------------*/
type PendingTask struct{
	TimeStamp time.Time
	Task db.UserRequest
}

var PendingChannel = make(chan PendingTask, 100);
var processingTime = 15*time.Second;

type PendingTaskQueue struct{
	Elements []PendingTask
}

func (pq *PendingTaskQueue) AddPendingTask(pendingTask PendingTask){
	currentTime := time.Now();
	timeDiff := currentTime.Sub(pendingTask.TimeStamp);

	if timeDiff < processingTime{
		waitTime := processingTime - timeDiff
		time.Sleep(waitTime);
	}

	res,err:= db.CheckPending(pendingTask.Task.Id)
	if err != nil{
		log.Fatalf("Failed to check the pending task %v", err);
	}

	if res {
		TaskChannel <- pendingTask.Task
		fmt.Print("Added the task for reprocessing")
	}
}