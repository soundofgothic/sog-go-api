package rest

import (
	"net/http"

	"github.com/ggicci/httpin"
	"soundofgothic.pl/backend/internal/domain"
	"soundofgothic.pl/backend/internal/rest/middlewares"
	"soundofgothic.pl/backend/internal/rest/rjson"
)

type SourceFilesListInput struct {
	Filter   string `in:"query=filter"`
	Page     int64  `in:"query=page;default=1" validate:"min=1"`
	PageSize int64  `in:"query=pageSize;default=50" validate:"min=10,max=100"`
	Type     string `in:"query=type"`

	GameIDs  middlewares.IDArray `in:"query=gameID,gameID[]"`
	NPCIDs   middlewares.IDArray `in:"query=npcID,npcID[]"`
	GuildIDs middlewares.IDArray `in:"query=guildID,guildID[]"`
	VoiceIDs middlewares.IDArray `in:"query=voiceID,voiceID[]"`
	IDs      middlewares.IDArray `in:"query=id,id[]"`
}

func (rc *restEnvironment) sourcefilesList(w http.ResponseWriter, r *http.Request) {
	input := r.Context().Value(httpin.Input).(*SourceFilesListInput)
	sourcefiles, count, err := rc.repositories.SourceFile().List(r.Context(), domain.SourceFileSearchOptions{
		Query:    input.Filter,
		Page:     input.Page,
		PageSize: input.PageSize,
		Type:     input.Type,

		GameIDs:  input.GameIDs.Values,
		GuildIDs: input.GuildIDs.Values,
		NPCIDs:   input.NPCIDs.Values,
		VoiceIDs: input.VoiceIDs.Values,
		IDs:      input.IDs.Values,
	})
	if err != nil {
		rjson.HandleError(w, err)
		return
	}
	rjson.RespondWithJSON(w, rjson.NewPagedResponse(rjson.PageOptions{Page: input.Page, PageSize: input.PageSize}, count, sourcefiles))
}
