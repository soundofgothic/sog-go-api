package rest

import (
	"net/http"

	"soundofgothic.pl/backend/ports/rest/middlewares"
	"soundofgothic.pl/backend/ports/rest/rjson"
)

func (re *restEnvironment) recordingsList(w http.ResponseWriter, r *http.Request) {
	page, err := middlewares.GetPageOptions(r.Context())
	if err != nil {
		rjson.HandleError(w, err)
		return
	}
	filter := r.URL.Query().Get("filter")
	recordings, total, err := re.repositories.Recording().SearchByText(r.Context(), filter, page.Page, page.PageSize)
	if err != nil {
		rjson.HandleError(w, err)
		return
	}
	rjson.RespondWithJSON(w, rjson.NewPagedResponse(rjson.PageOptions(*page), total, recordings))
}
