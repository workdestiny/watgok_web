package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/acoshift/csrf"
	"github.com/acoshift/httprouter"
	"github.com/acoshift/middleware"
	"github.com/acoshift/webstatic"
	"github.com/go-redis/redis"
	"github.com/moonrhythm/hime"
	"github.com/moonrhythm/session"
	"github.com/workdestiny/watgok_web/service"
)

// App is struct config for app
type App struct {
	SessionConfig session.Config
	Domain        string
	SQL           *sql.DB
	Redis         *redis.Client
	RedisPrefix   string
	CSRFConfig    csrf.Config
	Hime          *hime.App
	Location      *time.Location
}

// Config is the csrf config
type Config struct {
	Origins          []string
	ForbiddenHandler http.Handler
	IgnoreProto      bool
}

var (
	domain      string
	db          *sql.DB
	myRedis     *redis.Client
	redisPrefix string
	loc         *time.Location
	csef        csrf.Config
)

//Handler return Handler Muti Middleware
func (app *App) Handler() http.Handler {
	initConfig(app)

	mux := http.NewServeMux()
	mux.Handle("/healthz", http.HandlerFunc(healthzDatabaseHandler))
	mux.Handle("/-/", http.StripPrefix("/-", webstatic.New(webstatic.Config{
		Dir:          "public",
		CacheControl: "public, max-age=31536000",
	})))

	mux.Handle("/favicon.ico", fileHandler("public/images/favicon.ico"))

	// new http rounter
	m := httprouter.New()
	m.HandleMethodNotAllowed = false
	m.NotFound = hime.Handler(notFoundHandler)

	m.Get(app.Hime.Route("signin"), hime.Handler(signinGetHandle))

	// add m to mux
	mux.Handle("/", m)
	return middleware.Chain(
		DefaultCacheControl,
		logHTTP,
		noCORS,
		securityHeaders,
		methodFilter,
		csrf.New(app.CSRFConfig),
		session.Middleware(app.SessionConfig),
		panicRecovery,
	)(mux)
}

func initConfig(c *App) {
	domain = c.Domain
	myRedis = c.Redis
	redisPrefix = c.RedisPrefix
	db = c.SQL
	loc = c.Location
	csef = c.CSRFConfig
}

func healthzDatabaseHandler(w http.ResponseWriter, r *http.Request) {
	if err := db.Ping(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := myRedis.Do("PING").Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func notFoundHandler(ctx *hime.Context) error {
	return ctx.View("notfound", page(ctx))
}

func fileHandler(name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, name)
	})
}

func panicRecovery(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
				debug.PrintStack()

				// send massage to Discord Channel Chat
				service.SendErrorToDiscord(fmt.Sprintf("%s || %v", r.RequestURI, err))

				// w.Header().Set("Content-Type", "text/html; charset=utf-8")
				// w.Header().Set("X-Content-Type-Options", "nosniff")
				// w.Header().Set("Cache-Control", "private, no-cache, no-store, max-age=0")
				// w.WriteHeader(http.StatusInternalServerError)
				// io.Copy(w, bytes.NewReader(errorPage))
			}
		}()
		h.ServeHTTP(w, r)
	})
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
