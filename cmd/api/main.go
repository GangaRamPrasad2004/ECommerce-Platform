package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/GangaRamPrasad2004/learning-go-shop/internal/config"
	"github.com/GangaRamPrasad2004/learning-go-shop/internal/database"
	"github.com/GangaRamPrasad2004/learning-go-shop/internal/interfaces"
	"github.com/GangaRamPrasad2004/learning-go-shop/internal/logger"
	providers "github.com/GangaRamPrasad2004/learning-go-shop/internal/provider"
	"github.com/GangaRamPrasad2004/learning-go-shop/internal/server"
	"github.com/GangaRamPrasad2004/learning-go-shop/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	log := logger.New()
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	db, err := database.New(&cfg.Database)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect the database")
	}
	mainDB, err := db.DB()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get database connecttion")
	}
	defer func() {
		if err := mainDB.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close database connection")
		}
	}()
	gin.SetMode(cfg.Server.GinMode)

	authService := service.NewAuthService(db, cfg)
	userService := service.NewUserService(db)
	productService := service.NewProductService(db)

	var uploadProvider interfaces.UploadProvider
	if cfg.Upload.UploadProvider == "s3" {
		uploadProvider = providers.NewS3Provider(cfg)
	} else {
		uploadProvider = providers.NewLocalUploadProvider(cfg.Upload.Path)
	}
	uploadService := service.NewUploadService(uploadProvider)
	cartService := service.NewCartService(db)

	srv := server.New(cfg, db, &log, authService, userService, productService, uploadService, cartService)
	router := srv.SetupRoutes()
	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	go func() {
		log.Info().Str("port", cfg.Server.Port).Msg("starting http server")
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("failed to start http server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("failed to shutdown http server")
		return
	}
	log.Info().Msg("shutting down database")

}
