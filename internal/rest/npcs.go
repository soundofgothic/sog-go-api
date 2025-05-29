package rest

import (
	"net/http"

	"github.com/ggicci/httpin"
	"soundofgothic.pl/backend/internal/domain"
	"soundofgothic.pl/backend/internal/rest/middlewares"
	"soundofgothic.pl/backend/internal/rest/rjson"
)

type NPCListInput struct {
	Page      int64               `in:"query=page;default=1" validate:"min=1"`
	PageSize  int64               `in:"query=pageSize;default=50" validate:"min=10,max=100"`
	Filter    string              `in:"query=filter"`
	GameID    middlewares.IDArray `in:"query=gameID,gameID[]"`
	VoiceID   middlewares.IDArray `in:"query=voiceID,voiceID[]"`
	GuildID   middlewares.IDArray `in:"query=guildID,guildID[]"`
	ScriptIDs middlewares.IDArray `in:"query=scriptID,scriptID[]"`
	IDs       middlewares.IDArray `in:"query=id,id[]"`
}

func (re *restEnvironment) npcList(w http.ResponseWriter, r *http.Request) {
	input := r.Context().Value(httpin.Input).(*NPCListInput)
	npcs, count, err := re.repositories.NPC().List(r.Context(), domain.NPCSearchOptions{
		IDs:       input.IDs.Values,
		Query:     input.Filter,
		Page:      input.Page,
		PageSize:  input.PageSize,
		GameID:    input.GameID.Values,
		VoiceID:   input.VoiceID.Values,
		GuildID:   input.GuildID.Values,
		ScriptIDs: input.ScriptIDs.Values,
	})
	if err != nil {
		rjson.InternalError(w, err)
		return
	}
	rjson.RespondWithJSON(w, rjson.NewPagedResponse(rjson.PageOptions{Page: input.Page, PageSize: input.PageSize}, count, npcs))
}
