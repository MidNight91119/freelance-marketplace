package api

import (
	"fmt"
	"net/http"

	db "github.com/MidNight91119/freelance-marketplace/internal/db/sqlc"
	"github.com/MidNight91119/freelance-marketplace/internal/token"
	"github.com/MidNight91119/freelance-marketplace/internal/util"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	config     util.Config
	store      *db.Queries
	tokenMaker token.Maker
	validate   *validator.Validate
	router     *http.ServeMux
}

func NewServer(config util.Config, store *db.Queries) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
		validate:   validator.New(),
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/auth/signup", server.signup)
	mux.HandleFunc("POST /api/auth/login", server.login)

	server.router = mux
}

func (server *Server) Start(address string) error {
	return http.ListenAndServe(address, server.router)
}
