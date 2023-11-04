package rest

import (
	"net/http"

	"github.com/ggicci/httpin"
	"soundofgothic.pl/backend/domain"
	"soundofgothic.pl/backend/ports/rest/rjson"
)

type SourceFilesListInput struct {
	Filter   string `in:"query=filter"`
	GameID   int64  `in:"query=gameID"`
	Page     int64  `in:"query=page;default=1" validate:"min=1"`
	PageSize int64  `in:"query=pageSize;default=50" validate:"min=10,max=100"`
	Type     string `in:"query=type"`

	NPCIDs   []int64 `in:"query=npcID,npcID[]"`
	GuildIDs []int64 `in:"query=guildID,guildID[]"`
	VoiceIDs []int64 `in:"query=voiceID,voiceID[]"`
}

func (rc *restEnvironment) sourcefilesList(w http.ResponseWriter, r *http.Request) {
	input := r.Context().Value(httpin.Input).(*SourceFilesListInput)
	sourcefiles, count, err := rc.repositories.SourceFile().List(r.Context(), domain.SourceFileSearchOptions{
		Query:    input.Filter,
		GameID:   input.GameID,
		Page:     input.Page,
		PageSize: input.PageSize,
		Type:     input.Type,
	})
	if err != nil {
		rjson.HandleError(w, err)
		return
	}
	rjson.RespondWithJSON(w, rjson.NewPagedResponse(rjson.PageOptions{Page: input.Page, PageSize: input.PageSize}, count, sourcefiles))
}
