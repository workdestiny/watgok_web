package app

import (
	"log"

	"github.com/moonrhythm/hime"
)

func indexGetHandler(ctx *hime.Context) error {
	log.Println(GetMyID(ctx))
	return ctx.View("index", page(ctx))
}
