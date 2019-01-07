package main

import (
	"flag"
	"fmt"
	_ "reflect"

	"github.com/alexkylt/message_broker/internal/pkg/storage/mapstorage"

	"github.com/alexkylt/message_broker/internal/pkg/server"
	"github.com/alexkylt/message_broker/internal/pkg/storage/dbstorage"
)

func main() {
	var port int
	var mode string
	//var strg interface{}
	flag.IntVar(&port, "port", 9090, "specify port to use.  defaults to 9090.")
	flag.StringVar(&mode, "mode", "map", "specify storage to use.  defaults to map")
	flag.Parse()
	fmt.Println("MODE:", mode)
	switch mode {
	case "db":
		strg := dbstorage.InitStorage()
		srv := server.InitServer(port, strg)
		srv.Run()
	default:
		strg := mapstorage.InitStorage()
		srv := server.InitServer(port, strg)
		srv.Run()
	}
}
