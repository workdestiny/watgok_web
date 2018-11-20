package app

import (
	"github.com/moonrhythm/hime"
)

func postReadGetHandler(ctx *hime.Context) error {
	return ctx.View("app/post.read", page(ctx))
}
