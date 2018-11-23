package app

import (
	"log"

	"github.com/moonrhythm/hime"
	"github.com/workdestiny/watgok_web/config"
)

func signinGetHandle(ctx *hime.Context) error {

	log.Println("aloha!!!")

	return ctx.View("app/signin", page(ctx))
}

func signInFacebookGetHandler(ctx *hime.Context) error {
	return ctx.Redirect("https://www.facebook.com/v2.9/dialog/oauth",
		ctx.Param("client_id", config.FacebookAppID), ctx.Param("redirect_uri", baseURL+config.FacebookCallbackURL))
}
