package main

import (
	"log"
	"net"
	"net/rpc"
)

type Master struct {
	R        *rpc.Server
	TaskList chan TaskRes
}

var Ms Master

func (m *Master) run() {
	m.server()
}

func (m *Master) server() {
	t := new(TaskRes)
	m.TaskList = make(chan TaskRes, 10)
	m.R = rpc.NewServer()
	_ = m.R.Register(t)
	//_ = m.r.Register(m)
	l, e := net.Listen("tcp", ":9999")
	if e != nil {
		log.Fatal("rpc server start failed: ", e)
	}
	go func() {
		for {
			conn, err := l.Accept()
			if err == nil {
				go m.R.ServeConn(conn)
			} else {
				break
			}
		}
		_ = l.Close()
	}()
}
