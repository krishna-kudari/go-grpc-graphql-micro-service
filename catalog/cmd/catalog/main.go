package main

import (
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/krishna-kudari/go-grpc-graphql-micro-service/catalog"
	"github.com/tinrab/retry"
)

type Config struct {
	esURL string `envconfig:"ES_URL"`
}

func main() {
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatal(err)
	}
	var repository catalog.Repository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		repository, err = catalog.NewElasticRepository(config.esURL)
		if err != nil {
			log.Fatal(err)
		}
		return
	})
	defer repository.Close()
	log.Println("Listening on port 8080...")
	service := catalog.NewService(repository)
	log.Fatal(catalog.ListenGRPC(service, 8080))
}
