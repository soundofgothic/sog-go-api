package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"soundofgothic.pl/backend/internal/config"
	"soundofgothic.pl/backend/internal/domain"
	"soundofgothic.pl/backend/internal/postgres/migrations"
)

func testRepository() *postgresRepositoryStorage {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	repositories, err := NewPostgresRepositories(WithAuth(
		DBAuth{
			Host:     cfg.Postgres.Host,
			Port:     cfg.Postgres.Port,
			Username: cfg.Postgres.User,
			Password: cfg.Postgres.Password,
			Name:     cfg.Postgres.Database,
		},
	))
	if err != nil {
		panic(err)
	}

	if err := migrations.Migrate(context.Background(), repositories); err != nil {
		panic(err)
	}

	return repositories
}

func cleanupTestData(t *testing.T, r *postgresRepositoryStorage) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := r.db.NewDelete().Model((*domain.Alternative)(nil)).Where("1=1").Exec(ctx)
	require.NoError(t, err)
	_, err = r.db.NewDelete().Model((*domain.Mod)(nil)).Where("1=1").Exec(ctx)
	require.NoError(t, err)
	_, err = r.db.NewDelete().Model((*domain.TTSVoice)(nil)).Where("1=1").Exec(ctx)
	require.NoError(t, err)
}

func TestAlternativeRepository_List(t *testing.T) {
	r := testRepository()

	t.Cleanup(func() { cleanupTestData(t, r) })

	require.NoError(t, r.Mod().Persist(t.Context(), &domain.Mod{
		ID:          1,
		Name:        "Test Mod",
		Description: "Test Mod Description",
		Version:     "1.0.0",
		GameID:      2,
	}))

	require.NoError(t, r.TTSVoice().Persist(t.Context(), &domain.TTSVoice{
		ID:   1,
		Name: "Test Voice",
	}))

	require.NoError(t, r.Alternative().Persist(t.Context(), &domain.Alternative{
		ModID:       1,
		RecordingId: 1,
		TTSVoiceId:  1,
		Transcript:  "Test Transcript",
		State:       domain.AlternativeStateInProgress,
	}))

	require.NoError(t, r.Alternative().Persist(t.Context(), &domain.Alternative{
		ModID:       1,
		RecordingId: 2,
		TTSVoiceId:  1,
		Transcript:  "Test Transcript 2",
		State:       domain.AlternativeStateInProgress,
	}))

	t.Run("List with no filters", func(t *testing.T) {
		items, _, err := r.Alternative().List(t.Context(), domain.AlternativeOptions{})
		require.NoError(t, err)
		require.Len(t, items, 2)
		for _, item := range items {
			assert.NotNil(t, item.Recording)
			assert.NotNil(t, item.Recording.Game)
			assert.NotNil(t, item.Recording.Guild)
			assert.NotNil(t, item.Recording.Voice)
			assert.NotNil(t, item.Recording.SourceFile)
		}
	})

	testCases := []struct {
		name     string
		options  domain.AlternativeOptions
		expected int
	}{
		{
			name: "List with modID filter - all",
			options: domain.AlternativeOptions{
				ModID: 1,
			},
			expected: 2,
		},
		{
			name: "list with modID filter - none",
			options: domain.AlternativeOptions{
				ModID: 2,
			},
			expected: 0,
		},
		{
			name: "list with modID and recordingID filter - one",
			options: domain.AlternativeOptions{
				ModID: 1,
				RecordingSearchOptions: domain.RecordingSearchOptions{
					IDs: []int64{1},
				},
			},
			expected: 1,
		},
		{
			name: "list with string query on recording transcript - one",
			options: domain.AlternativeOptions{
				RecordingSearchOptions: domain.RecordingSearchOptions{
					Query: "zostać członkiem",
				},
			},
			expected: 1,
		},
		{
			name: "list with string query on alternative transcript - one",
			options: domain.AlternativeOptions{
				Query: "Test Transcript 2",
			},
			expected: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			items, _, err := r.Alternative().List(t.Context(), tc.options)
			require.NoError(t, err)
			require.Len(t, items, tc.expected)
		})
	}
}

func TestAlternativeRepository_Persist(t *testing.T) {
	r := testRepository()

	t.Cleanup(func() { cleanupTestData(t, r) })

	// Setup dependencies
	require.NoError(t, r.Mod().Persist(t.Context(), &domain.Mod{
		ID:          1,
		Name:        "Test Mod",
		Description: "Test Mod Description",
		Version:     "1.0.0",
		GameID:      2,
	}))
	require.NoError(t, r.TTSVoice().Persist(t.Context(), &domain.TTSVoice{
		ID:   1,
		Name: "Test Voice",
	}))

	// Test create new Alternative
	alt := &domain.Alternative{
		ModID:       1,
		RecordingId: 10,
		TTSVoiceId:  1,
		Transcript:  "Initial Transcript",
		State:       domain.AlternativeStateInProgress,
	}
	err := r.Alternative().Persist(t.Context(), alt)
	require.NoError(t, err)

	// Fetch and verify
	items, _, err := r.Alternative().List(t.Context(), domain.AlternativeOptions{ModID: 1, RecordingSearchOptions: domain.RecordingSearchOptions{IDs: []int64{10}}})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "Initial Transcript", items[0].Transcript)

	// Test update Alternative
	alt.Transcript = "Updated Transcript"
	alt.State = domain.AlternativeStateCompleted
	err = r.Alternative().Persist(t.Context(), alt)
	require.NoError(t, err)

	// Fetch and verify update
	items, _, err = r.Alternative().List(t.Context(), domain.AlternativeOptions{ModID: 1, RecordingSearchOptions: domain.RecordingSearchOptions{IDs: []int64{10}}})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "Updated Transcript", items[0].Transcript)
	assert.Equal(t, domain.AlternativeStateCompleted, items[0].State)
}

func TestAlternativeRepository_Delete(t *testing.T) {
	r := testRepository()
	t.Cleanup(func() { cleanupTestData(t, r) })

	// Setup dependencies
	require.NoError(t, r.Mod().Persist(t.Context(), &domain.Mod{
		ID:          1,
		Name:        "Test Mod",
		Description: "Test Mod Description",
		Version:     "1.0.0",
		GameID:      2,
	}))
	require.NoError(t, r.TTSVoice().Persist(t.Context(), &domain.TTSVoice{
		ID:   1,
		Name: "Test Voice",
	}))

	// Create and persist an alternative
	alt := &domain.Alternative{
		ModID:       1,
		RecordingId: 20,
		TTSVoiceId:  1,
		Transcript:  "To be deleted",
		State:       domain.AlternativeStateInProgress,
	}
	require.NoError(t, r.Alternative().Persist(t.Context(), alt))

	// Ensure it exists
	items, _, err := r.Alternative().List(t.Context(), domain.AlternativeOptions{ModID: 1, RecordingSearchOptions: domain.RecordingSearchOptions{IDs: []int64{20}}})
	require.NoError(t, err)
	require.Len(t, items, 1)

	// Delete the alternative
	err = r.Alternative().Delete(t.Context(), alt.ModID, alt.RecordingId)
	require.NoError(t, err)

	// Ensure it no longer exists
	items, _, err = r.Alternative().List(t.Context(), domain.AlternativeOptions{ModID: 1, RecordingSearchOptions: domain.RecordingSearchOptions{IDs: []int64{20}}})
	require.NoError(t, err)
	require.Len(t, items, 0)
}
