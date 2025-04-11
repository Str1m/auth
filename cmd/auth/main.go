package main

import (
	"context"
	"github.com/Str1m/auth/internal/app"
	"log"
)

const DSN = "test"

func main() {

	//cfg, err := env.NewPGConfig()
	//if err != nil {
	//	log.Fatalf("failed to get pg config: %s", err.Error())
	//}
	//

	//
	//p, err := pgxpool.New(ctx, cfg.DSN())
	//if err != nil {
	//	log.Fatalln("err")
	//}
	//
	//client := postgres.NewClientPG(p)
	//
	//DBLayer := postgres.NewStoragePG(client)
	ctx := context.Background()

	a, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("failed to init app %s", err.Error())
	}

	err = a.Run()
	if err != nil {
		log.Fatalf("failed to run app %s", err.Error())
	}
}
