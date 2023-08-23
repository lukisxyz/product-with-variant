package main

import (
	"context"
	"flukis/product/cmd"
	"flukis/product/config"
	"flukis/product/internals/attribute"
	"flukis/product/internals/category"
	"flukis/product/internals/product"
	"flukis/product/internals/product_category"
	"flukis/product/internals/variant"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
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
	attributeRepo := attribute.NewRepo(pool)
	attributeSvc := attribute.NewService(
		attributeRepo,
		pool,
	)
	attributeRouter := attribute.NewRouter(attributeSvc)

	// attr
	categoryRepo := category.NewRepo(pool)
	categorySvc := category.NewService(
		categoryRepo,
		pool,
	)
	categoryRouter := category.NewRouter(categorySvc)

	// attr
	productVariantRepo := variant.NewRepo(pool)
	productVariantSvc := variant.NewService(
		productVariantRepo,
		pool,
	)
	productVariantRouter := variant.NewRouter(productVariantSvc)

	// attr
	productCategory := product_category.NewRepo(pool)
	productRepo := product.NewRepo(pool)
	productSvc := product.NewService(
		productRepo,
		productCategory,
		pool,
	)
	productRouter := product.NewRouter(productSvc)

	// Create router.
	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Mount("/attribute", attributeRouter.Routes())
	r.Mount("/category", categoryRouter.Routes())
	r.Mount("/product", productRouter.Routes())
	r.Mount("/variant", productVariantRouter.Routes())

	// Run server instance.
	log.Info().Msg("starting up server...")
	if err := cmd.Run(&cfg, r); err != nil {
		log.Fatal().Err(err).Msg("failed to start the server")
		return
	}
	log.Info().Msg("server Stopped")
}
