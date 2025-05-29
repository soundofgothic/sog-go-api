package rest

import (
	"net/http"

	"soundofgothic.pl/backend/internal/rest/rjson"
)

func (re *restEnvironment) gamesList(w http.ResponseWriter, r *http.Request) {
	games, err := re.repositories.Game().List(r.Context())
	if err != nil {
		rjson.InternalError(w, err)
		return
	}
	rjson.RespondWithJSON(w, games)
}
