package rest

import (
	"net/http"

	"soundofgothic.pl/backend/ports/rest/rjson"
)

func (re *restEnvironment) voicesList(w http.ResponseWriter, r *http.Request) {
	voices, err := re.repositories.Voice().List(r.Context())
	if err != nil {
		rjson.HandleError(w, err)
		return
	}
	rjson.RespondWithJSON(w, voices)
}
