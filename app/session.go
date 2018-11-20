package app

import (
	"github.com/moonrhythm/hime"
	"github.com/moonrhythm/session"
)

var (
	// SessionName is sess
	SessionName = "sess"
	// SessionUserKey is key in session
	SessionUserKey = "userID"
)

func getSession(ctx *hime.Context) *session.Session {
	s, err := session.Get(ctx, SessionName)
	must(err)
	return s
}

func setSession(ctx *hime.Context, userID string) error {
	s, err := session.Get(ctx, SessionName)
	if err != nil {
		return err
	}
	s.Set(SessionUserKey, userID)
	return err
}

func removeSession(ctx *hime.Context) error {
	s, err := session.Get(ctx, SessionName)
	if err != nil {
		return err
	}
	s.Del(SessionUserKey)
	return nil

}
