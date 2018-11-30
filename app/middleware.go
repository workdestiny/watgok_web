package app

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/acoshift/header"
	"github.com/acoshift/middleware"
	"github.com/workdestiny/watgok_web/entity"
	"github.com/workdestiny/watgok_web/repository"
)

// DefaultCacheControl sets default cache-control header
func DefaultCacheControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(header.CacheControl, "no-cache, no-store, must-revalidate")
		h.ServeHTTP(w, r)
	})
}

func logHTTP(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timeStart := time.Now().In(loc)
		lrw := newLoggingResponseWriter(w)
		h.ServeHTTP(lrw, r)
		if strings.HasPrefix(r.URL.Path, "/-") {
			return
		}
		timeEnd := time.Now().In(loc)
		clientIP, port, _ := net.SplitHostPort(r.RemoteAddr)
		statusCode := lrw.statusCode
		request := r.RequestURI
		fmt.Printf("%s | %v | %v | %3d | %13v | %s | %s\n", timeEnd.Format("2006-01-02 15:04:05"), clientIP, port, statusCode, timeEnd.Sub(timeStart), r.Method, request)
	})
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func noCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func securityHeaders(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(header.XFrameOptions, "deny")
		w.Header().Set(header.XXSSProtection, "1; mode=block")
		w.Header().Set(header.XContentTypeOptions, "nosniff")
		h.ServeHTTP(w, r)
	})
}

func methodFilter(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet, http.MethodHead, http.MethodPost, http.MethodOptions:
			h.ServeHTTP(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
}

func getCookie(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID := getUserID(ctx)
		nctx := WithMyID(ctx, userID)
		h.ServeHTTP(w, r.WithContext(nctx))
	})
}

func fetchUser(db *sql.DB) middleware.Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rCtx := r.Context()
			u := &entity.UserModel{}
			id := GetMyID(rCtx)
			if id != "" {
				u, _ = repository.GetUser(db, id)
				if u.ID == "" {
					removeSession(r.Context())
				}
			}

			ctx := WithUser(rCtx, u)
			ctx = WithUserRole(ctx, entity.Role(u.Role))
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
