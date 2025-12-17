package main

import (
	"log"
	"os"

	"github.com/easyhire/backend/internal/executor"
	"github.com/gin-gonic/gin"
)

func main() {
	addr := os.Getenv("EXECUTOR_ADDR")
	if addr == "" {
		addr = "0.0.0.0:8090"
	}

	runner := executor.NewRunner()

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	srv := executor.NewHTTPServer(runner)
	srv.Register(router)

	log.Printf("🚀 EasyHire Executor listening on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatal(err)
	}
}
