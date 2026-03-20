package survey

import (
	"testing"
	"time"
)

func TestCalculatePhaseBoundaries(t *testing.T) {
	createdAt := time.Date(2026, 3, 20, 10, 0, 0, 0, time.UTC)
	voteEndsAt := createdAt.Add(24 * time.Hour)
	resultsEndsAt := createdAt.Add(48 * time.Hour)

	tests := []struct {
		name string
		now  time.Time
		want Phase
	}{
		{
			name: "before vote end is VOTING",
			now:  voteEndsAt.Add(-time.Nanosecond),
			want: PhaseVoting,
		},
		{
			name: "at vote end is RESULTS",
			now:  voteEndsAt,
			want: PhaseResults,
		},
		{
			name: "before results end is RESULTS",
			now:  resultsEndsAt.Add(-time.Nanosecond),
			want: PhaseResults,
		},
		{
			name: "at results end is EXPIRED",
			now:  resultsEndsAt,
			want: PhaseExpired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculatePhase(tt.now, voteEndsAt, resultsEndsAt)
			if got != tt.want {
				t.Fatalf("CalculatePhase() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestComputeFlags(t *testing.T) {
	base := time.Date(2026, 3, 20, 10, 0, 0, 0, time.UTC)
	voteEndsAt := base.Add(24 * time.Hour)
	resultsEndsAt := base.Add(48 * time.Hour)

	t.Run("private pin open live during voting", func(t *testing.T) {
		got := ComputeFlags(
			base,
			VisibilityPrivatePIN,
			ResultsModeOpenLive,
			voteEndsAt,
			resultsEndsAt,
		)

		if got.Phase != PhaseVoting {
			t.Fatalf("Phase = %s, want %s", got.Phase, PhaseVoting)
		}
		if !got.CanVote {
			t.Fatalf("CanVote = false, want true")
		}
		if !got.ResultsVisible {
			t.Fatalf("ResultsVisible = false, want true")
		}
		if !got.RequiresPIN {
			t.Fatalf("RequiresPIN = false, want true")
		}
	})

	t.Run("closed hidden until end at vote end", func(t *testing.T) {
		got := ComputeFlags(
			voteEndsAt,
			VisibilityPublic,
			ResultsModeClosedHiddenUntilEnd,
			voteEndsAt,
			resultsEndsAt,
		)

		if got.Phase != PhaseResults {
			t.Fatalf("Phase = %s, want %s", got.Phase, PhaseResults)
		}
		if got.CanVote {
			t.Fatalf("CanVote = true, want false")
		}
		if !got.ResultsVisible {
			t.Fatalf("ResultsVisible = false, want true")
		}
		if got.RequiresPIN {
			t.Fatalf("RequiresPIN = true, want false")
		}
	})

	t.Run("closed hidden until end before vote end", func(t *testing.T) {
		got := ComputeFlags(
			voteEndsAt.Add(-time.Nanosecond),
			VisibilityUnlisted,
			ResultsModeClosedHiddenUntilEnd,
			voteEndsAt,
			resultsEndsAt,
		)

		if got.Phase != PhaseVoting {
			t.Fatalf("Phase = %s, want %s", got.Phase, PhaseVoting)
		}
		if !got.CanVote {
			t.Fatalf("CanVote = false, want true")
		}
		if got.ResultsVisible {
			t.Fatalf("ResultsVisible = true, want false")
		}
		if got.RequiresPIN {
			t.Fatalf("RequiresPIN = true, want false")
		}
	})
}
