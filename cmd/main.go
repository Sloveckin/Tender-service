package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"net/http"
	"tender-service/internal/config"
	"tender-service/internal/initialization"
)

func main() {
	cnf := config.MustLoadConfig()
	logger := initialization.InitLogger()
	repos := initialization.RepositoriesInit(cnf)
	services := initialization.ServiceInit(repos)
	handlers := initialization.HandlersInit(services, logger)

	server := initialization.InitChiServer(handlers)

	logger.Info("Service Run")
	err := http.ListenAndServe(":"+cnf.ServerAddress, server)
	if err != nil {
		panic(fmt.Sprintf("Not expected error: %s", err))
	}
}
