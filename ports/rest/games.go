package rest

import (
	"net/http"

	"soundofgothic.pl/backend/ports/rest/rjson"
)

func (re *restEnvironment) gamesList(w http.ResponseWriter, r *http.Request) {
	games, err := re.repositories.Game().List(r.Context())
	if err != nil {
		rjson.HandleError(w, err)
		return
	}
	rjson.RespondWithJSON(w, games)
}
