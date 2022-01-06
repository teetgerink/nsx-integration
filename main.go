package main

import (
	"os"

	"github.com/alecthomas/template"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/teetgerink/nsx-integration.git/api"
)

// This is here to keep the github.com/alecthomas/template in the go.mod, without it the build fails.
var _ = template.HTMLEscape

func main() {
	r := gin.New()
	r.Use()

	// Determine where to run
	httpServerAddress := os.Getenv("HTTP_SERVER_ADDRESS")
	if httpServerAddress == "" {
		httpServerAddress = "0.0.0.0:8080"
	}

	// Configure our application
	nsxapi := api.NSXAPI{}
	runner := nsxapi.Configure(r)

	logrus.Info("Starting webserver")
	logrus.Fatal(runner.Run(httpServerAddress))
}
