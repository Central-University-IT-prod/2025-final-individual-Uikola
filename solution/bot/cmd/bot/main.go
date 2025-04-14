package main

import (
	"bot/internal/app"
	"context"

	"github.com/rs/zerolog/log"
)

func main() {
	ctx := context.Background()

	a, err := app.NewApp(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to init app")
		return
	}

	a.Run()
}
