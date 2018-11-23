package app

import (
	"encoding/json"
	"image"
	"log"
	"net/http"

	"github.com/moonrhythm/hime"
	"github.com/workdestiny/watgok_web/config"
	"github.com/workdestiny/watgok_web/entity"
	"github.com/workdestiny/watgok_web/service"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
)

func signinGetHandle(ctx *hime.Context) error {

	log.Println("aloha!!!")

	return ctx.View("app/signin", page(ctx))
}

func signInFacebookGetHandler(ctx *hime.Context) error {
	return ctx.Redirect("https://www.facebook.com/v2.9/dialog/oauth",
		ctx.Param("client_id", config.FacebookAppID),
		ctx.Param("redirect_uri", baseURL+config.FacebookCallbackURL))
}

func signInFacebookCallbackGetHandler(ctx *hime.Context) error {

	facebookOauth2 := oauth2.Config{
		ClientID:     config.FacebookAppID,
		ClientSecret: fbToken,
		RedirectURL:  baseURL + config.FacebookCallbackURL,
		Scopes:       []string{"public_profile", "email"},
		Endpoint:     facebook.Endpoint,
	}

	code := ctx.Request.URL.Query().Get("code")

	//https: //graph.facebook.com/v2.9/oauth/access_token?client_id=appid&redirect_uri=link&client_secret&code=code
	tokenFacebook, err := facebookOauth2.Exchange(ctx, code)
	if err != nil {
		return err
	}

	//https: //graph.facebook.com/v2.9/me?fields=id,name,email,picture.type(large)&access_token=
	resp, err := http.Get("https://graph.facebook.com/v2.9/me?fields=id,name,email,picture.type(large)&access_token=" + tokenFacebook.AccessToken)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var facebookData entity.FacebookOauth2
	err = json.NewDecoder(resp.Body).Decode(&facebookData)
	if err != nil {
		return err
	}

	resp, err = http.Get(facebookData.Picture.Data.URL)
	if err != nil {
		return ctx.RedirectTo("signin")
	}
	defer resp.Body.Close()
	image, _, err := image.Decode(resp.Body)

	displayImage := service.ResizeDisplay(image)
	path := service.GenerateDisplayName(facebookData.ID)
	upload(ctx, displayImage, path)

	return nil
}
