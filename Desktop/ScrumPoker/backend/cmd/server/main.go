package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"

	"scrum-poker/backend/internal/cache"
	"scrum-poker/backend/internal/config"
	httplayer "scrum-poker/backend/internal/http"
	"scrum-poker/backend/internal/service"
	"scrum-poker/backend/internal/store"
	"scrum-poker/backend/internal/ws"
)

func main() {
	cfg := config.Load()

	db, err := store.OpenMySQL(cfg.MySQLDSN())
	if err != nil {
		log.Fatalf("open mysql: %v", err)
	}
	defer db.Close()

	if err := store.RunMigrations(db); err != nil {
		log.Fatalf("run migrations: %v", err)
	}

	redisCache, err := cache.NewRedisCache(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB, cfg.PresenceTTL)
	if err != nil {
		log.Fatalf("open redis: %v", err)
	}
	defer redisCache.Close()

	repo := store.NewRepository(db)
	roomService := service.NewRoomService(repo, redisCache)
	hub := ws.NewHub(redisCache, roomService)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go hub.Run(ctx)

	app := fiber.New(fiber.Config{
		AppName:       "Scrum Poker MVP",
		CaseSensitive: true,
	})

	httplayer.Register(app, httplayer.Dependencies{
		Config:  cfg,
		Service: roomService,
		Hub:     hub,
	})

	go func() {
		if err := app.Listen(":" + cfg.Port); err != nil {
			log.Printf("fiber stopped: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop

	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		log.Printf("shutdown fiber: %v", err)
	}
}
