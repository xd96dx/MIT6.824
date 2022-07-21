package main

import (
	"fmt"
	"log"
	"net/rpc"
)

type TaskRes struct {
	Id       int
	Res      bool
	TaskName string
}

func connect() *rpc.Client {
	client, err := rpc.Dial("tcp", ":9999")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	return client
}

func (t *TaskRes) Put(task *TaskRes, reply *string) error {
	fmt.Println("rpc put: ", task)
	if !task.Res {
		return nil
	}
	Ms.TaskList <- *task
	fmt.Println("len: ", len(Ms.TaskList))
	*reply = "success"
	return nil
}

func (t *TaskRes) Get(reply *TaskRes, noUse *string) error {
	if len(Ms.TaskList) == 0 {
		return nil
	}
	*reply = <- Ms.TaskList
	fmt.Println("reply: ", reply)
	return nil
}