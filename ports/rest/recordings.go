package rest

import (
	"net/http"

	"github.com/ggicci/httpin"
	"soundofgothic.pl/backend/domain"
	"soundofgothic.pl/backend/ports/rest/rjson"
)

type RecordingListInput struct {
	Page         int64  `in:"query=page;default=1"`
	PageSize     int64  `in:"query=pageSize;default=50"`
	Filter       string `in:"query=filter"`
	SourceFileID int64  `in:"query=sourceFileID"`
	VoiceID      int64  `in:"query=voiceID"`
	GameID       int64  `in:"query=gameID"`
	NPCID        int64  `in:"query=npcID"`
}

func (re *restEnvironment) recordingsList(w http.ResponseWriter, r *http.Request) {
	input := r.Context().Value(httpin.Input).(*RecordingListInput)
	recordings, total, err := re.repositories.Recording().List(r.Context(), domain.RecordingSearchOptions{
		Query:        input.Filter,
		Page:         input.Page,
		PageSize:     input.PageSize,
		GameID:       input.GameID,
		SourceFileID: input.SourceFileID,
		NPCID:        input.NPCID,
		VoiceID:      input.VoiceID,
	})
	if err != nil {
		rjson.HandleError(w, err)
		return
	}
	rjson.RespondWithJSON(w, rjson.NewPagedResponse(rjson.PageOptions{Page: input.Page, PageSize: input.PageSize}, total, recordings))
}
