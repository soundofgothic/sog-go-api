package domain

import (
	"context"

	"github.com/uptrace/bun"
)

type AlternativeState string

const (
	AlternativeStateNone       AlternativeState = "none"
	AlternativeStateScheduled  AlternativeState = "scheduled"
	AlternativeStateInProgress AlternativeState = "in_progress"
	AlternativeStateCompleted  AlternativeState = "completed"
	AlternativeStateFailed     AlternativeState = "failed"
)

type Alternative struct {
	bun.BaseModel `bun:"alternatives,alias:a"`
	ModID         int64      `json:"modID" bun:"mod_id,notnull,pk"`
	Mod           *Mod       `json:"mod,omitempty" bun:"rel:belongs-to,join:mod_id=id"`
	RecordingId   int64      `json:"recordingID" bun:"recording_id,notnull,pk"`
	Recording     *Recording `json:"recording,omitempty" bun:"r,rel:belongs-to,join:recording_id=id"`
	TTSVoiceId    int64      `json:"ttsVoiceID" bun:"tts_voice_id,notnull"`
	TTSVoice      *TTSVoice  `json:"ttsVoice,omitempty" bun:"rel:belongs-to,join:tts_voice_id=id"`
	Transcript    string     `json:"transcript" bun:"transcript,notnull"`

	State AlternativeState `json:"state" bun:"state,notnull,default:none"`
	Wave  string           `json:"wave" bun:"wave"`
}

type AlternativeOptions struct {
	RecordingSearchOptions
	ModID      int64  `search:"type:exact;columns:a.mod_id;"`
	Query      string `search:"type:like;columns:a.transcript;"`
	State      string `search:"type:exact;columns:a.state;"`
	TTSVoiceID int64  `search:"type:exact;columns:a.tts_voice_id;"`
}

type AlternativeService interface {
	Persist(ctx context.Context, alternative *Alternative) error
	List(ctx context.Context, query AlternativeOptions) ([]Alternative, int64, error)
	Get(ctx context.Context, modID, recordingID int64) (*Alternative, error)
	Delete(ctx context.Context, modID, recordingID int64) error
}
