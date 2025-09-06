package main

//
//  Swagger API docs
//  @title           Music Analyze API
//  @version         1.0
//  @description     API endpoints for music analytics derived from Chinook dataset.
//  @BasePath        /

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/go-chi/chi/v5"
    "github.com/bryanwahyu/analyze/internal/config"
    appanalyst "github.com/bryanwahyu/analyze/internal/application/analyst"
    analyticsHTTP "github.com/bryanwahyu/analyze/internal/http/analytics"
    mysqlp "github.com/bryanwahyu/analyze/internal/infra/db/mysql"
    httpSwagger "github.com/swaggo/http-swagger"
    docs "github.com/bryanwahyu/analyze/docs"
)

func main() {
	// path config.yaml
	path := "config.yaml"
	if v := os.Getenv("CONFIG_PATH"); v != "" {
		path = v
	}

	// load config
	cfg, err := config.Load(path)
	if err != nil {
		log.Fatalf("config load error: %v", err)
	}

	ctx := context.Background()

	// connect MySQL
	db, err := mysqlp.Connect(ctx, cfg.MySQLDSN())
	if err != nil {
		log.Fatalf("mysql connect error: %v", err)
	}
	defer db.Close()

    // init repo and service
    analystRepo := mysqlp.NewAnalystRepository(db)
    analystSvc := appanalyst.NewService(analystRepo)

    // init router
    mux := chi.NewRouter()

    // swagger docs metadata (override at runtime if needed)
    docs.SwaggerInfo.Title = "Music Analyze API"
    docs.SwaggerInfo.Description = "API endpoints for music analytics derived from Chinook dataset."
    docs.SwaggerInfo.Version = "1.0"
    docs.SwaggerInfo.BasePath = "/"

    // swagger UI endpoint
    mux.Get("/swagger/*", httpSwagger.Handler(
        httpSwagger.URL("/swagger/doc.json"),
    ))

    // app routes
    mux.Mount("/", analyticsHTTP.NewRouter(analystSvc))

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// run server
	go func() {
		log.Printf("server listening on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	log.Println("shutting down server...")

	ctx2, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx2); err != nil {
		log.Printf("shutdown error: %v", err)
	}
}
