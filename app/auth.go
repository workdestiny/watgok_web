package app

import (
	"log"

	"github.com/moonrhythm/hime"
)

func signinGetHandle(ctx *hime.Context) error {

	log.Println("aloha!!!")

	return ctx.View("app/signin", page(ctx))
}
