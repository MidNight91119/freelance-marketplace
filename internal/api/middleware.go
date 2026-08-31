package api

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const (
	authTypeBearer            = "bearer"
	authHeaderKey             = "authorization"
	authPayloadKey contextKey = "auth_payload"
)

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
