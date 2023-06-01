package rest

import (
	"net/http"

	"github.com/ggicci/httpin"
	"soundofgothic.pl/backend/domain"
	"soundofgothic.pl/backend/ports/rest/rjson"
)

type GuildListInput struct {
	Page        int64  `in:"query=page;default=1" validate:"min=1"`
	PageSize    int64  `in:"query=pageSize;default=50" validate:"min=10,max=100"`
	GameID      int64  `in:"query=gameID"`
	InGameAlias string `in:"query=ingameAlias"`
	Filter      string `in:"query=filter"`
}

func (re *restEnvironment) guildsList(w http.ResponseWriter, r *http.Request) {
	input := r.Context().Value(httpin.Input).(*GuildListInput)
	guilds, count, err := re.repositories.Guild().List(r.Context(), domain.GuildSearchOptions{
		Query:       input.Filter,
		Page:        input.Page,
		PageSize:    input.PageSize,
		GameID:      input.GameID,
		InGameAlias: input.InGameAlias,
	})
	if err != nil {
		rjson.HandleError(w, err)
		return
	}
	rjson.RespondWithJSON(w, rjson.NewPagedResponse(rjson.PageOptions{Page: input.Page, PageSize: input.PageSize}, count, guilds))
}
