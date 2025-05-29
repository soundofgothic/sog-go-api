package rest

import (
	"net/http"

	"github.com/ggicci/httpin"
	"soundofgothic.pl/backend/internal/domain"
	"soundofgothic.pl/backend/internal/rest/middlewares"
	"soundofgothic.pl/backend/internal/rest/rjson"
)

type GuildListInput struct {
	Page        int64               `in:"query=page;default=1" validate:"min=1"`
	PageSize    int64               `in:"query=pageSize;default=50" validate:"min=10,max=100"`
	GameIDs     middlewares.IDArray `in:"query=gameID,gameID[]"`
	InGameAlias string              `in:"query=ingameAlias"`
	Filter      string              `in:"query=filter"`
	VoiceIDs    middlewares.IDArray `in:"query=voiceID,voiceID[]"`
	ScriptIDs   middlewares.IDArray `in:"query=scriptID,scriptID[]"`
	NPCIDs      middlewares.IDArray `in:"query=npcID,npcID[]"`
	IDs         middlewares.IDArray `in:"query=id,id[]"`
}

func (re *restEnvironment) guildsList(w http.ResponseWriter, r *http.Request) {
	input := r.Context().Value(httpin.Input).(*GuildListInput)
	guilds, count, err := re.repositories.Guild().List(r.Context(), domain.GuildSearchOptions{
		Query:       input.Filter,
		Page:        input.Page,
		PageSize:    input.PageSize,
		GameIDs:     input.GameIDs.Values,
		InGameAlias: input.InGameAlias,
		VoiceIDs:    input.VoiceIDs.Values,
		ScriptIDs:   input.ScriptIDs.Values,
		NPCIDs:      input.NPCIDs.Values,
		IDs:         input.IDs.Values,
	})
	if err != nil {
		rjson.InternalError(w, err)
		return
	}
	rjson.RespondWithJSON(w, rjson.NewPagedResponse(rjson.PageOptions{Page: input.Page, PageSize: input.PageSize}, count, guilds))
}
