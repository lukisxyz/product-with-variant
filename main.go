package main

import (
	"context"
	"flukis/product/cmd"
	"flukis/product/config"
	"flukis/product/internals/product_attributes"
	"flukis/product/internals/product_categories"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func main() {
	// load config
	err := godotenv.Load()
	cfg := config.LoadConfig("config.yml")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load file .env")
	}

	// database
	dbString := cfg.DBConfig.ConnString()
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbString)
	if err != nil {
		log.Error().Err(err).Msg("unable to connect to database")
	}
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load file .env")
	}

	// attr
	attributeRepo := product_attributes.NewRepo(pool)
	attributeSvc := product_attributes.NewService(
		attributeRepo,
		pool,
	)
	attributeRouter := product_attributes.NewRouter(attributeSvc)

	// attr
	categoryRepo := product_categories.NewRepo(pool)
	categorySvc := product_categories.NewService(
		categoryRepo,
		pool,
	)
	categoryRouter := product_categories.NewRouter(categorySvc)

	// Create router.
	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)

	r.Mount("/attribute", attributeRouter.Routes())
	r.Mount("/category", categoryRouter.Routes())

	// Run server instance.
	log.Info().Msg("starting up server...")
	if err := cmd.Run(&cfg, r); err != nil {
		log.Fatal().Err(err).Msg("failed to start the server")
		return
	}
	log.Info().Msg("server Stopped")
}
