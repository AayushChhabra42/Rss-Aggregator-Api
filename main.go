package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/AayushChhabra42/rssaggregator/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {

	godotenv.Load(".env")

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("Port not Configured")
	}
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not found")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("connection error", err)
	}
	db := database.New(conn)
	apiCfg := apiConfig{
		DB: db,
	}
	go startScraping(db, 10, time.Minute)
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1router := chi.NewRouter()
	v1router.Get("/healthz", handlerReadiness)
	v1router.Get("/err", handlerErr)

	v1router.Post("/users", apiCfg.handlerCreateUser)
	v1router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerGetUser))

	v1router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerFeedCreate))
	v1router.Get("/feeds", apiCfg.handlerGetFeeds)

	v1router.Get("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerFeedFollowsGet))
	v1router.Post("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerFeedFollowCreate))
	v1router.Delete("/feed_follows/{feedFollowID}", apiCfg.middlewareAuth(apiCfg.handlerFeedFollowDelete))

	router.Mount("/v1", v1router)

	server := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}
	log.Printf("Server starting on port %v", portString)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
