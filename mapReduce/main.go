package main

import (
	"fmt"
	"time"
)

func putTask() {
	c := connect()
	n := 1
	for {
		name := fmt.Sprintf("task-%d", n)
		task := TaskRes{
			Id:       n,
			Res:      true,
			TaskName: name,
		}
		var res string
		err := c.Call("TaskRes.Put", task, &res)
		if err != nil {
			fmt.Println(err)
			return
		}
		time.Sleep(time.Second)
		n++
	}
}

func getTask() {
	c := connect()
	for {
		task := new(TaskRes)
		var s string
		err := c.Call("TaskRes.Get", task, &s)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func rpcRun() {
	go putTask()
	go getTask()
}

func main() {
	Ms.run()
	rpcRun()
	time.Sleep(3 * time.Second)
}