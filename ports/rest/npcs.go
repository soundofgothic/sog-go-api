package rest

import (
	"net/http"

	"github.com/ggicci/httpin"
	"soundofgothic.pl/backend/domain"
	"soundofgothic.pl/backend/ports/rest/rjson"
)

type NPCListInput struct {
	Page     int64  `in:"query=page;default=1"`
	PageSize int64  `in:"query=pageSize;default=50"`
	Filter   string `in:"query=filter"`
	GameID   int64  `in:"query=gameID"`
	VoiceID  int64  `in:"query=voiceID"`
	GuildID  int64  `in:"query=guildID"`
}

func (re *restEnvironment) npcList(w http.ResponseWriter, r *http.Request) {
	input := r.Context().Value(httpin.Input).(*NPCListInput)
	npcs, count, err := re.repositories.NPC().List(r.Context(), domain.NPCSearchOptions{
		Query:    input.Filter,
		Page:     input.Page,
		PageSize: input.PageSize,
		GameID:   input.GameID,
		VoiceID:  input.VoiceID,
		GuildID:  input.GuildID,
	})
	if err != nil {
		rjson.HandleError(w, err)
		return
	}
	rjson.RespondWithJSON(w, rjson.NewPagedResponse(rjson.PageOptions{Page: input.Page, PageSize: input.PageSize}, count, npcs))
}
