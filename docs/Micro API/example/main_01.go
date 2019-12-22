package main

import (
	"github.com/micro/micro/cmd"

	// go-micro plugins
	_ "github.com/micro/go-plugins/registry/consul"
	_ "github.com/micro/go-plugins/registry/kubernetes"
	_ "github.com/micro/go-plugins/transport/grpc"
	_ "github.com/micro/go-plugins/transport/tcp"
)

func main() {
	cmd.Init()
}
