package main

import (
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/krishna-kudari/go-grpc-graphql-micro-service/account"
	"github.com/tinrab/retry"
)

type Config struct {
	databaseURL string `envconfig:"DATABASE_URL"`
}

func main()  {
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatal(err)
	}

	var r account.Repository
	retry.ForeverSleep(2*time.Second, func(_ int)(err error){
		r, err = account.NewPostgresRepository(config.databaseURL)
		if err != nil {
			log.Println(err)
		}
		return
	})
	defer r.Close()
	log.Println("Listening on port 8080...")
	s := account.NewService(r)
	log.Fatal(account.ListenGRPC(s, 8080))
}
