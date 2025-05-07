package rest

import (
	"net/http"

	"github.com/ggicci/httpin"
	"soundofgothic.pl/backend/internal/domain"
	"soundofgothic.pl/backend/internal/rest/rjson"
)

type RecordingListInput struct {
	Page         int64   `in:"query=page;default=1" validate:"min=1"`
	PageSize     int64   `in:"query=pageSize;default=50" validate:"min=10,max=100"`
	Filter       string  `in:"query=filter"`
	SourceFileID []int64 `in:"query=sourceFileID,sourceFileID[]"`
	VoiceID      []int64 `in:"query=voiceID,voiceID[]"`
	GameID       []int64 `in:"query=gameID,gameID[]"`
	NPCID        []int64 `in:"query=npcID,npcID[]"`
	GuildID      []int64 `in:"query=guildID,guildID[]"`
}

func (re *restEnvironment) recordingsList(w http.ResponseWriter, r *http.Request) {
	input := r.Context().Value(httpin.Input).(*RecordingListInput)
	recordings, total, err := re.repositories.Recording().List(r.Context(), domain.RecordingSearchOptions{
		Query:        input.Filter,
		Page:         input.Page,
		PageSize:     input.PageSize,
		GameIDs:      input.GameID,
		SourceFileID: input.SourceFileID,
		NPCIDs:       input.NPCID,
		VoiceIDs:     input.VoiceID,
		GuildIDs:     input.GuildID,
	})
	if err != nil {
		rjson.HandleError(w, err)
		return
	}
	rjson.RespondWithJSON(w, rjson.NewPagedResponse(rjson.PageOptions{Page: input.Page, PageSize: input.PageSize}, total, recordings))
}

type RecordingGetInput struct {
	GameID int64  `in:"path=gameID" validate:"required,min=1,max=2"`
	Wave   string `in:"path=wave" validate:"required"`
}

func (re *restEnvironment) recordingsGet(w http.ResponseWriter, r *http.Request) {
	input := r.Context().Value(httpin.Input).(*RecordingGetInput)
	record, err := re.repositories.Recording().Get(r.Context(), domain.RecordingGetOptions{
		GameID: input.GameID,
		Wave:   input.Wave,
	})
	if err != nil {
		rjson.HandleError(w, err)
		return
	}
	rjson.RespondWithJSON(w, record)
}
