package matcher

type Result struct {
	Input      string
	URL        string
	StatusCode int
	Size       int
	Words      int
	Lines      int
}

type Rules struct {
	MatchStatus  map[int]bool
	FilterStatus map[int]bool

	MatchSize  map[int]bool
	FilterSize map[int]bool

	MatchWords  map[int]bool
	FilterWords map[int]bool

	MatchLines  map[int]bool
	FilterLines map[int]bool
}

func Matches(r Result, rules Rules) bool {
	// Match conditions.
	if len(rules.MatchStatus) > 0 &&
		!rules.MatchStatus[r.StatusCode] {
		return false
	}

	if len(rules.MatchSize) > 0 &&
		!rules.MatchSize[r.Size] {
		return false
	}

	if len(rules.MatchWords) > 0 &&
		!rules.MatchWords[r.Words] {
		return false
	}

	if len(rules.MatchLines) > 0 &&
		!rules.MatchLines[r.Lines] {
		return false
	}

	// Filter conditions.
	if rules.FilterStatus[r.StatusCode] {
		return false
	}

	if rules.FilterSize[r.Size] {
		return false
	}

	if rules.FilterWords[r.Words] {
		return false
	}

	if rules.FilterLines[r.Lines] {
		return false
	}

	return true
}
