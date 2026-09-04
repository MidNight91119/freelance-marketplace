package api

import (
	"context"
	"net/http"
	"slices"
	"strings"

	"github.com/MidNight91119/freelance-marketplace/internal/token"
)

type contextKey string

const (
	authTypeBearer            = "bearer"
	authHeaderKey             = "authorization"
	authPayloadKey contextKey = "auth_payload"
)

func payloadFrom(r *http.Request) (*token.Payload, bool) {
	payload, ok := r.Context().Value(authPayloadKey).(*token.Payload)
	return payload, ok
}

func (server *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get(authHeaderKey)
		if len(authHeader) == 0 {
			writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "authorization header is not provided")
			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) < 2 {
			writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "invalid authorization header format")
			return
		}

		if strings.ToLower(fields[0]) != authTypeBearer {
			writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "unsupported authorization scheme")
			return
		}
		accessToken := fields[1]

		payload, err := server.tokenMaker.VerifyToken(accessToken)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "invalid or expired access token")
			return
		}

		ctx := context.WithValue(r.Context(), authPayloadKey, payload)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (server *Server) requireRole(next http.Handler, allowed ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload, ok := payloadFrom(r)
		if !ok {
			writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "something went wrong")
			return
		}

		if !slices.Contains(allowed, payload.Role) {
			writeError(w, http.StatusForbidden, "FORBIDDEN", "you are not allowed to perform this action")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (server *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
