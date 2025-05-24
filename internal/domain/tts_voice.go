package domain

import "context"

type TTSVoice struct {
	ID int64 `json:"id" bun:"id,pk"`
	// Name is the name of the TTS voice.
	Name string `json:"name" bun:"name,notnull"`
}

type TTSVoiceService interface {
	Persist(ctx context.Context, voice *TTSVoice) error
	List(ctx context.Context) ([]TTSVoice, error)
	Get(ctx context.Context, id int64) (*TTSVoice, error)
	Delete(ctx context.Context, id int64) error
}
