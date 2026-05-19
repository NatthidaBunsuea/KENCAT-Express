// โอม
package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"kencatexpress/backend/internal/config"
	"kencatexpress/backend/internal/util"
)

type contextKey string

const claimsKey contextKey = "auth_claims"

type Middleware func(http.Handler) http.Handler

func Chain(h http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

func Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				util.ErrorJSON(w, http.StatusInternalServerError, "internal server error")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		recorder := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(recorder, r)
		log.Printf("%s %s %d %s", r.Method, r.URL.Path, recorder.status, time.Since(start))
	})
}

func CORS(cfg config.Config) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			allowedOrigin := strings.TrimSpace(cfg.CORSAllowOrigin)
			requestOrigin := strings.TrimSpace(r.Header.Get("Origin"))

			// อนุญาตเฉพาะ origin ที่กำหนดไว้เท่านั้น
			if requestOrigin != "" && requestOrigin != allowedOrigin {
				util.ErrorJSON(w, http.StatusForbidden, "origin not allowed")
				return
			}

			if allowedOrigin != "" {
				w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
				w.Header().Set("Vary", "Origin")
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")
			}

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func OptionalAuth(secret string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := readClaims(r, secret)
			if ok {
				r = r.WithContext(context.WithValue(r.Context(), claimsKey, claims))
			}
			next.ServeHTTP(w, r)
		})
	}
}

func RequireAuth(secret string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := readClaims(r, secret)
			if !ok {
				util.ErrorJSON(w, http.StatusUnauthorized, "authentication required")
				return
			}
			r = r.WithContext(context.WithValue(r.Context(), claimsKey, claims))
			next.ServeHTTP(w, r)
		})
	}
}

func RequireRoles(roles ...string) Middleware {
	allowed := map[string]struct{}{}
	for _, role := range roles {
		allowed[util.NormalizeRole(role)] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := ClaimsFromContext(r.Context())
			if !ok {
				util.ErrorJSON(w, http.StatusUnauthorized, "authentication required")
				return
			}
			if _, exists := allowed[util.NormalizeRole(claims.Role)]; !exists {
				util.ErrorJSON(w, http.StatusForbidden, "insufficient permissions")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func ClaimsFromContext(ctx context.Context) (*util.Claims, bool) {
	claims, ok := ctx.Value(claimsKey).(*util.Claims)
	return claims, ok
}

func readClaims(r *http.Request, secret string) (*util.Claims, bool) {
	token := util.ParseBearerToken(r.Header.Get("Authorization"))
	if token == "" {
		cookie, err := r.Cookie(util.CookieName)
		if err == nil {
			token = strings.TrimSpace(cookie.Value)
		}
	}
	if token == "" {
		return nil, false
	}

	claims, err := util.ParseJWT(token, secret)
	if err != nil {
		return nil, false
	}
	return claims, true
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}
