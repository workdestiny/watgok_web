package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/acoshift/configfile"
	"github.com/acoshift/csrf"
	"github.com/acoshift/probehandler"
	"github.com/go-redis/redis"
	_ "github.com/lib/pq"
	"github.com/moonrhythm/hime"
	"github.com/moonrhythm/session"
	redisstore "github.com/moonrhythm/session/store/goredis"
	"github.com/workdestiny/watgok_web/app"
)

func main() {

	configValue := configfile.NewYAMLReader("config/application.yaml")

	loc, _ := time.LoadLocation("Asia/Bangkok")

	sessionHost := configValue.String("session_host")
	redisPrefix := configValue.String("session_prefix")
	redisClient := redis.NewClient(&redis.Options{
		Addr:       sessionHost,
		MaxRetries: 3,
		PoolSize:   6,
	})
	defer redisClient.Close()

	db, err := sql.Open("postgres", configValue.String("db"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.SetMaxIdleConns(2)
	db.SetMaxOpenConns(0)
	db.SetConnMaxLifetime(30 * time.Minute)
	db.Ping()

	appHime := hime.New()

	appFactory := &app.App{
		Location:    loc,
		Domain:      configValue.String("domain"),
		SQL:         db,
		Redis:       redisClient,
		RedisPrefix: redisPrefix,
		SessionConfig: session.Config{
			Store: redisstore.New(redisstore.Config{
				Client: redisClient,
				Prefix: redisPrefix,
			}),
			HTTPOnly: true,
			Secure:   session.PreferSecure,
			Proxy:    true,
			MaxAge:   60 * 24 * time.Hour,
			Path:     "/",
			Rolling:  true,
			Keys:     [][]byte{configValue.Bytes("session_key")},
			Secret:   configValue.Bytes("session_secret"),
			SameSite: http.SameSiteLaxMode,
		},
		CSRFConfig: csrf.Config{
			Origins: []string{
				configValue.String("domain"),
				"watgok.local:8080",
			},
			IgnoreProto: true,
		},
		Hime:   appHime,
		Static: static("public/mix-manifest.json"),
	}
	appHime.Template().
		Funcs(appFactory.TemplateFuncs()).
		ParseConfigFile("settings/web/template.yaml")
	appHime.
		ParseConfigFile("settings/web/routes.yaml").
		ParseConfigFile("settings/web/server.yaml").
		Handler(appFactory.Handler())

	probe := probehandler.New()
	health := http.NewServeMux()
	health.Handle("/", probehandler.Success())
	health.Handle("/readiness", probe)
	go http.ListenAndServe(":18080", health)

	appHime.
		GracefulShutdown().
		Notify(probe.Fail)

	log.Println("Web Running!")

	err = appHime.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func static(filename string) map[string]string {
	s := make(map[string]string)
	bs, _ := ioutil.ReadFile(filename)
	json.Unmarshal(bs, &s)
	return s
}
