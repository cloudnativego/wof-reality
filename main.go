package main

import (
	"fmt"
	"os"

	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/cloudnativego/wof-reality/service"
)

func main() {
	appEnv, err := cfenv.Current()
	if err != nil {
		fmt.Printf("CF environment not detected. APP WILL RUN WITH FAKE REPOSITORY!\n")
	}

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3000"
	}

	server := service.NewServer(appEnv)
	server.Run(":" + port)
}
