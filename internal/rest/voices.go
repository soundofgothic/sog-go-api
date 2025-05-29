package rest

import (
	"net/http"

	"github.com/ggicci/httpin"
	"soundofgothic.pl/backend/internal/domain"
	"soundofgothic.pl/backend/internal/rest/middlewares"
	"soundofgothic.pl/backend/internal/rest/rjson"
)

type VoicesListInput struct {
	Filter    string              `in:"query=filter"`
	GameIDs   middlewares.IDArray `in:"query=gameID,gameID[]"`
	NPCIDs    middlewares.IDArray `in:"query=npcID,npcID[]"`
	GuildIDs  middlewares.IDArray `in:"query=guildID,guildID[]"`
	ScriptIDs middlewares.IDArray `in:"query=scriptID,scriptID[]"`
	IDs       middlewares.IDArray `in:"query=id,id[]"`
}

func (re *restEnvironment) voicesList(w http.ResponseWriter, r *http.Request) {
	input := r.Context().Value(httpin.Input).(*VoicesListInput)
	voices, err := re.repositories.Voice().List(r.Context(), domain.VoiceOptions{
		Query:     input.Filter,
		GuildIDs:  input.GuildIDs.Values,
		NPCIDs:    input.NPCIDs.Values,
		GameIDs:   input.GameIDs.Values,
		ScriptIDs: input.ScriptIDs.Values,
	})
	if err != nil {
		rjson.InternalError(w, err)
		return
	}
	rjson.RespondWithJSON(w, voices)
}
