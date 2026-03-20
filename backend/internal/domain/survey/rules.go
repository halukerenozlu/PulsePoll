package survey

import "time"

type Phase string

const (
	PhaseVoting  Phase = "VOTING"
	PhaseResults Phase = "RESULTS"
	PhaseExpired Phase = "EXPIRED"
)

type Visibility string

const (
	VisibilityPublic     Visibility = "public"
	VisibilityUnlisted   Visibility = "unlisted"
	VisibilityPrivatePIN Visibility = "private_pin"
)

type ResultsMode string

const (
	ResultsModeOpenLive             ResultsMode = "open_live"
	ResultsModeClosedHiddenUntilEnd ResultsMode = "closed_hidden_until_end"
)

type ComputedFlags struct {
	Phase          Phase
	CanVote        bool
	ResultsVisible bool
	RequiresPIN    bool
}

func CalculatePhase(now, voteEndsAt, resultsEndsAt time.Time) Phase {
	if now.Before(voteEndsAt) {
		return PhaseVoting
	}
	if now.Before(resultsEndsAt) {
		return PhaseResults
	}
	return PhaseExpired
}

func ComputeFlags(
	now time.Time,
	visibility Visibility,
	resultsMode ResultsMode,
	voteEndsAt time.Time,
	resultsEndsAt time.Time,
) ComputedFlags {
	phase := CalculatePhase(now, voteEndsAt, resultsEndsAt)

	resultsVisible := false
	switch resultsMode {
	case ResultsModeOpenLive:
		resultsVisible = true
	case ResultsModeClosedHiddenUntilEnd:
		resultsVisible = phase != PhaseVoting
	}

	return ComputedFlags{
		Phase:          phase,
		CanVote:        phase == PhaseVoting,
		ResultsVisible: resultsVisible,
		RequiresPIN:    visibility == VisibilityPrivatePIN,
	}
}
