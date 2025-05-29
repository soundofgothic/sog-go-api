package postgres

import (
	"testing"

	"github.com/stretchr/testify/require"
	"soundofgothic.pl/backend/internal/domain"
)

func TestRecordingRepository(t *testing.T) {
	r := testRepository()

	t.Cleanup(func() { cleanupTestData(t, r) })

	testCases := []struct {
		name            string
		opts            domain.RecordingSearchOptions
		extepectedCount int
	}{
		{
			name: "dupa",
			opts: domain.RecordingSearchOptions{
				Query: "dupa",
			},
			extepectedCount: 13,
		},
		{
			name: "dupa with npcs",
			opts: domain.RecordingSearchOptions{
				Query:  "dupa",
				NPCIDs: []int64{814, 1006},
			},
			extepectedCount: 2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			records, count, err := r.Recording().List(t.Context(), tc.opts)
			require.NoError(t, err)
			require.Len(t, records, tc.extepectedCount)
			require.EqualValues(t, tc.extepectedCount, count)
		})
	}
}
