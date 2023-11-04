package rest

import (
	"net/http"

	"github.com/ggicci/httpin"
	"soundofgothic.pl/backend/domain"
	"soundofgothic.pl/backend/ports/rest/rjson"
)

type VoicesListInput struct {
	GameIDs   []int64 `in:"query=gameID,gameID[]"`
	NPCIDs    []int64 `in:"query=npcID,npcID[]"`
	GuildIDs  []int64 `in:"query=guildID,guildID[]"`
	ScriptIDs []int64 `in:"query=scriptID,scriptID[]"`
}

func (re *restEnvironment) voicesList(w http.ResponseWriter, r *http.Request) {
	input := r.Context().Value(httpin.Input).(*VoicesListInput)
	voices, err := re.repositories.Voice().List(r.Context(), domain.VoiceOptions{
		GuildIDs:  input.GuildIDs,
		NPCIDs:    input.NPCIDs,
		GameIDs:   input.GameIDs,
		ScriptIDs: input.ScriptIDs,
	})
	if err != nil {
		rjson.HandleError(w, err)
		return
	}
	rjson.RespondWithJSON(w, voices)
}
