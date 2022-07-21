package main

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

type ctl struct {
	funcType string
	keyMap   map[int]string
}

type worker struct {
	id     int
	result int
	task   string
	mutex  sync.Mutex
}

func writeFile(ctl *ctl) error {
	f, err := os.OpenFile("", os.O_RDWR | os.O_APPEND | os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	inData, _ := json.Marshal(ctl.keyMap)
	_, err = f.Write(inData)
	_ = f.Close()
	return err
}

func mapFunc(ctl *ctl) error {
	return nil
}

func reduceFunc() {

}

func (w *worker) worker(ctl ctl) {
	switch ctl.funcType {
	case "map":
		log.Printf("map func %d start", w.id)
		err := mapFunc(&ctl)
		// if map func failed, send the response to master to re-execute.
		if err != nil {
			c := connect()
			task := TaskRes{
				Id:       w.id,
				Res:      false,
				TaskName: w.task,
			}
			//var res string
			err = c.Call("TaskRes.Put", task, "")
		}
	case "reduce":
		reduceFunc()
	}
}
