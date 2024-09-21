package main

import (
	"fmt"
	"net/rpc"
	"os"

	"github.com/fatih/color"
)

func rpcClient(inMaintenance bool) error {
	c, err := rpc.Dial("tcp", "127.0.0.1:"+os.Getenv("RPC_PORT"))
	if err != nil {
		return err
	}

	fmt.Println("Connected...")
	var result string
	err = c.Call("RPCServer.MaintenanceMode", inMaintenance, &result)
	if err != nil {
		return err
	}

	color.Yellow(result)
	return nil
}
