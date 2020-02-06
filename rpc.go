package main

import (
	"fmt"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

// a single socket RPC listener
func listenAndServeRPC(listen net.Listener, server *rpc.Server, maximumConnections uint64) {
accept_loop:
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Printf("rpc.Server terminated: accept error:%v\n", err)
			break accept_loop
		}
		if connectionCountRPC.Increment() <= maximumConnections {
			go func() {
				server.ServeCodec(jsonrpc.NewServerCodec(conn))
				conn.Close()
				connectionCountRPC.Decrement()
			}()
		} else {
			connectionCountRPC.Decrement()
			conn.Close()
		}

	}
	listen.Close()
	fmt.Println("RPC accept terminated")
}
