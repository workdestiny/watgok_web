package service

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/acoshift/configfile"
	"github.com/workdestiny/watgok_web/entity"
)

//SendErrorToDiscord send err
func SendErrorToDiscord(err string) {
	url := configfile.NewYAMLReader("config/config-stage.yaml").String("webhook")
	if url == "" {
		return
	}

	dc := entity.Discord{
		Content: err,
	}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(dc)
	http.Post(url, "application/json; charset=utf-8", b)
}
