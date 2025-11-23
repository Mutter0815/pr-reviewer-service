package main

import (
	"log"

	"github.com/Mutter0815/pr-reviewer-service/internal/app"
	httptransport "github.com/Mutter0815/pr-reviewer-service/internal/transport/http"
)

func main() {
	application := app.New()
	defer application.Close()

	router := httptransport.NewRouter(application.Services)
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("failed to start server:%v", err)
	}

}
