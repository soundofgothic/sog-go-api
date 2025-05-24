package domain

type Repositories interface {
	Game() GameService
	Guild() GuildService
	Voice() VoiceService
	SourceFile() SourceFileService
	Recording() RecordingService
	NPC() NPCService

	Mod() ModService
	Alternative() AlternativeService
	TTSVoice() TTSVoiceService
}
