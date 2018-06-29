package main

import (
	"context"
	"time"

	"github.com/ariefrahmansyah/spinnaker-demo/web"
)

func main() {
	ctxWeb := context.Background()

	serverOptions := &web.Options{
		ListenAddress:  ":8080",
		MaxConnections: 512,
		ReadTimeout:    10 * time.Second,
	}

	webServer := web.New(nil, serverOptions)

	if err := webServer.Run(ctxWeb); err != nil {
		panic(err)
	}
}
