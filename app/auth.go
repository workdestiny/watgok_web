package app

import (
	"database/sql"
	"encoding/json"
	"image"
	"net/http"

	"github.com/acoshift/pgsql"
	"github.com/moonrhythm/hime"
	"github.com/workdestiny/watgok_web/config"
	"github.com/workdestiny/watgok_web/entity"
	"github.com/workdestiny/watgok_web/repository"
	"github.com/workdestiny/watgok_web/service"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
)

func signinGetHandle(ctx *hime.Context) error {
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

	u, err := repository.GetUserByFBID(db, facebookData.ID)
	if err == sql.ErrNoRows {
		//create
		resp, err = http.Get(facebookData.Picture.Data.URL)
		if err != nil {
			return ctx.RedirectTo("index")
		}
		defer resp.Body.Close()

		image, _, err := image.Decode(resp.Body)
		if err != nil {
			return err
		}

		displayImage := service.ResizeDisplay(image)
		path := service.GenerateDisplayName(facebookData.ID)
		upload(ctx, displayImage, path)
		//url image display
		url := generateDownloadURL(path)

		err = pgsql.RunInTx(db, nil, func(tx *sql.Tx) error {
			//insert
			id, err := repository.CreateUser(db, facebookData.Name, url, facebookData.ID, entity.RoleUser)
			if err != nil {
				return err
			}

			//signin
			setSession(ctx, id)
			return nil
		})
		if err != nil {
			return err
		}

		return ctx.RedirectTo("index")
	}
	must(err)
	//signin
	setSession(ctx, u.ID)
	return ctx.RedirectTo("index")
}
