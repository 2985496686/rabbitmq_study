package rpc

import (
	"testing"
)

func TestRpcClient(t *testing.T) {
	for i := 0; i < 5; i++ {
		rpcClient()
	}
}

func TestRpcServer(t *testing.T) {
	rpcServer()
}
