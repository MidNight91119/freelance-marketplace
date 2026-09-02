package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	db "github.com/MidNight91119/freelance-marketplace/internal/db/sqlc"
	"github.com/MidNight91119/freelance-marketplace/internal/token"
	"github.com/MidNight91119/freelance-marketplace/internal/util"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type signupRequest struct {
	Name     string `json:"name" validate:"required,max=255"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=6"`
	Role     string `json:"role" validate:"required,oneof=client freelancer"`
}

type userResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      string(user.Role),
		CreatedAt: user.CreatedAt,
	}
}

type loginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type loginResponse struct {
	AccessToken          string       `json:"accessToken"`
	AccessTokenExpiresAt time.Time    `json:"accessTokenExp"`
	User                 userResponse `json:"user"`
}

func (server *Server) signup(w http.ResponseWriter, r *http.Request) {
	var req signupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "request body is not valid JSON")
		return
	}

	err := server.validate.Struct(req)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		log.Printf("signup: hash password: %v", err)
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "something went wrong")
		return
	}

	arg := db.CreateUserParams{
		Name:           req.Name,
		Email:          req.Email,
		HashedPassword: hashedPassword,
		Role:           db.Roles(req.Role),
	}

	user, err := server.store.CreateUser(r.Context(), arg)
	if err != nil {
		pgErr, ok := errors.AsType[*pgconn.PgError](err)
		if ok && pgErr.Code == db.UniqueViolation {
			writeError(w, http.StatusConflict, "EMAIL_ALREADY_EXISTS", "an account with this email already exists")
			return
		}
		log.Printf("signup: create user failed: %v", err)
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "something went wrong")
		return
	}

	writeJSON(w, http.StatusCreated, newUserResponse(user))
}

func (server *Server) login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "request body is not valid JSON")
		return
	}

	if err := server.validate.Struct(req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	user, err := server.store.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeError(w, http.StatusUnauthorized, "INVALID_CREDENTIALS", "invalid email or password")
			return
		}
		log.Printf("login: get user failed: %v", err)
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "something went wrong")
		return
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "INVALID_CREDENTIALS", "invalid email or password")
		return
	}

	accessToken, payload, err := server.tokenMaker.CreateToken(user.ID, user.Email, string(user.Role), server.config.AccessTokenDuration)
	if err != nil {
		log.Printf("login: create token failed: %v", err)
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "something went wrong")
		return
	}

	rsp := loginResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: payload.ExpiredAt,
		User:                 newUserResponse(user),
	}

	writeJSON(w, http.StatusOK, rsp)
}

func (server *Server) getMe(w http.ResponseWriter, r *http.Request) {
	payload, ok := r.Context().Value(authPayloadKey).(*token.Payload)
	if !ok {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "middleware not working")
		return
	}

	user, err := server.store.GetUserByEmail(r.Context(), payload.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeError(w, http.StatusUnauthorized, "INVALID_CREDENTIALS", "invalid email or password")
			return
		}
		log.Printf("getMe: get user failed: %v", err)
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "something went wrong")
		return
	}

	writeJSON(w, http.StatusOK, newUserResponse(user))
}
