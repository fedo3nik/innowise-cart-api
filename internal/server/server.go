package server

import (
	"cart-api/internal/pkg/common/db/repository"
	"cart-api/internal/pkg/common/endpoints"
	"cart-api/internal/pkg/config"
	middleware "cart-api/internal/server/middleware"
	"log"
	"net/http"

	_ "cart-api/docs"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/jmoiron/sqlx"
)

type Server struct {
	cfg  *config.Config
	pool *sqlx.DB
}

func NewServer(config *config.Config, dbPool *sqlx.DB) *Server {
	return &Server{
		cfg:  config,
		pool: dbPool,
	}
}

func (s *Server) registerRoutes() *http.ServeMux {
	cartHandler := endpoints.NewCarHandler(s.pool)
	cartItemHandler := endpoints.NewCartItemHandler(s.pool)

	cartRepository := repository.NewPostgresCartRepository(s.pool)
	cartItemRepository := repository.NewPostgresItemRepository(s.pool)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /carts", cartHandler.CreateNew(cartRepository))
	mux.HandleFunc("GET /carts/{id}", cartHandler.ViewCart(cartRepository))
	mux.HandleFunc("GET /carts", cartHandler.GetAll(cartRepository))
	mux.HandleFunc("DELETE /carts/{id}", cartHandler.DeleteCart(cartRepository))

	wrappedPath := middleware.NewValiDateMiddleWare(cartItemHandler.AddToCart(cartItemRepository))

	mux.Handle("POST /carts/{cartId}/items", wrappedPath)
	mux.HandleFunc("DELETE /carts/{cartId}/items/{itemId}", cartItemHandler.RemoveFromCart(cartItemRepository, cartRepository))
	mux.Handle("GET /swagger/", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:3000/swagger/doc.json"),
	))
	return mux
}

func (s *Server) Run() error {

	router := s.registerRoutes()

	wrappedMux := middleware.NewLoggingMiddleware(router)

	muxServer := &http.Server{
		Addr:    s.cfg.GetPort(),
		Handler: wrappedMux,
	}

	log.Printf("Server listen on port:  %s", s.cfg.GetPort())

	if err := muxServer.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
