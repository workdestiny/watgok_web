package app

import (
	"github.com/moonrhythm/hime"
)

func indexGetHandler(ctx *hime.Context) error {

	return ctx.View("index", page(ctx))
}
