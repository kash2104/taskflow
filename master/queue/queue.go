package queue

import (
	"fmt"

	"github.com/kash2104/taskflow/master/db"
)


type Queue struct{
	Elements []db.UserRequest
}

func (q *Queue) Push(req db.UserRequest){
	q.Elements = append(q.Elements, req);
}

func (q *Queue) Pop() (error){
	if q.IsEmpty(){
		return fmt.Errorf("queue is empty");
	}

	q.Elements = q.Elements[1:]
	return nil;
}

func (q *Queue) Top() (db.UserRequest,error){
	if q.IsEmpty(){
		return db.UserRequest{},fmt.Errorf("queue is empty");
	}

	return q.Elements[0], nil
}

func (q *Queue) IsEmpty() bool{
	return len(q.Elements) == 0;
}
