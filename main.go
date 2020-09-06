package main

import (
	"log"
        "flag"
        "fmt"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"

	"lightbulb.org/bto"
        "lightbulb.org/defaults"
)

func routing(r *gin.Engine, ctrl *Controller) {
	r.GET("/on", ctrl.On)
	r.GET("/off", ctrl.Off)
	r.GET("/status", ctrl.Status)
	r.GET("/setBrightness", ctrl.SetBrightness)
        r.GET("/getBrightness", ctrl.GetBrightness)
}

func main() {
        var (
                host = flag.String("host", "localhost", "host address")
                port = flag.Int("port", 50051, "port number")
                irPath = flag.String("ir", "irdata.toml", "ir data path")
        )
        flag.Parse()

        var conf bto.Config
        irDataHander := defaults.New(*irPath)
        irDataHander.Load(&conf)

        address := fmt.Sprintf("%v:%d", *host, *port)
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := bto.NewIRServiceClient(conn)
        irclient := bto.NewLightbulbController(client, conf)
        ctrl := NewController(irclient)
	r := gin.Default()
	routing(r, ctrl)

	r.Run(":3001")
}
