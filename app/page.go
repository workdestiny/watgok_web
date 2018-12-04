package app

import (
	"net/http"
	"time"

	"github.com/moonrhythm/hime"
)

func page(ctx *hime.Context) map[string]interface{} {
	sess := getSession(ctx)
	r := ctx.Request

	x := make(map[string]interface{})
	x["QueryURI"] = r.URL.Query().Get("type")
	if r.URL.Path != "" {
		x["Path"] = r.URL.Path
	}
	x["Tagline"] = "| swapgap เขียนบทความ เขียนบล็อก อัพเดททุกวัน"
	x["URL"] = "https://" + getHost(r) + r.RequestURI
	x["Flash"] = sess.Flash()
	x["Now"] = time.Now()
	x["User"] = getUser(ctx)

	return x
}

// getHost gets real host from request
func getHost(r *http.Request) string {
	host := r.Header.Get("X-Forwarded-Host")
	if host == "" {
		host = r.Host
	}
	return host
}
