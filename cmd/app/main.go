package main

import (
	"context"
	"kafkaExperiments/internal/app"
	"log"
)

func main() {
	ctx := context.Background()

	app := app.New()

	if err := app.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
