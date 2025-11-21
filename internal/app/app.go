package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Vladimirmoscow84/Event_Booker/internal/handlers"
	"github.com/Vladimirmoscow84/Event_Booker/internal/notifier"
	"github.com/Vladimirmoscow84/Event_Booker/internal/service"
	"github.com/Vladimirmoscow84/Event_Booker/internal/storage"
	"github.com/Vladimirmoscow84/Event_Booker/internal/storage/postgres"
	"github.com/wb-go/wbf/config"
	"github.com/wb-go/wbf/ginext"
)

func Run() {

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	cfg := config.New()
	err := cfg.LoadEnvFiles(".env")
	if err != nil {
		log.Fatalf("[app] error of loading cfg: %v", err)
	}
	cfg.EnableEnv("")

	databaseURI := cfg.GetString("DATABASE_URI")

	serverAddr := cfg.GetString("SERVER_ADDRESS")

	eHost := cfg.GetString("EMAIL_HOST")
	ePort := cfg.GetString("EMAIL_PORT")
	eUser := cfg.GetString("EMAIL_USER")
	ePass := cfg.GetString("EMAIL_PASS")
	eFrom := cfg.GetString("EMAIL_FROM")
	eTo := cfg.GetString("EMAIL_TO")

	pgStore, err := postgres.New(databaseURI)
	if err != nil {
		log.Fatalf("[app] failed to connect to PG DB: %v", err)
	}
	defer pgStore.Close()

	storage, err := storage.New(pgStore)
	if err != nil {
		log.Fatalf("[app] failed to init unified storage: %v", err)
	}
	log.Println("[app] storage initialized successfully")

	emailClient := notifier.New(eHost, ePort, eUser, ePass, eFrom, eTo)
	log.Println("[app] email client initialized successfully")

	bookingTTL := 30 * time.Minute

	serviceClient := service.New(storage, emailClient, bookingTTL)
	log.Println("[app] service initialized successfully")

	serviceClient.StartBookingWorker(ctx, 1*time.Minute)

	engine := ginext.New("release")

	router := handlers.New(engine, serviceClient, serviceClient, serviceClient, serviceClient, serviceClient)
	router.Routes()

	srv := &http.Server{
		Addr:    serverAddr,
		Handler: engine,
	}

	// запуск сервера в отдельной горутине
	go func() {
		log.Printf("[app] starting server on %s", serverAddr)
		err = srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("[app] server failed: %v", err)
		}
	}()

	//ShutDown
	<-ctx.Done()
	log.Println("[app] shutting down server...")

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctxShutdown); err != nil {
		log.Fatalf("[app] server shutdown failed: %v", err)
	}

	log.Println("[app] server stopped")

}
