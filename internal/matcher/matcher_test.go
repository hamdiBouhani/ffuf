package matcher

import "testing"

func TestMatches(t *testing.T) {
	tests := []struct {
		name   string
		result Result
		rules  Rules
		want   bool
	}{
		{
			name: "matches status code",
			result: Result{
				StatusCode: 200,
				Size:       100,
				Words:      10,
				Lines:      5,
			},
			rules: Rules{
				MatchStatus: map[int]bool{
					200: true,
				},
			},
			want: true,
		},
		{
			name: "does not match status code",
			result: Result{
				StatusCode: 404,
			},
			rules: Rules{
				MatchStatus: map[int]bool{
					200: true,
				},
			},
			want: false,
		},
		{
			name: "filters status code",
			result: Result{
				StatusCode: 404,
			},
			rules: Rules{
				FilterStatus: map[int]bool{
					404: true,
				},
			},
			want: false,
		},
		{
			name: "matches size",
			result: Result{
				StatusCode: 200,
				Size:       1234,
			},
			rules: Rules{
				MatchSize: map[int]bool{
					1234: true,
				},
			},
			want: true,
		},
		{
			name: "filters size",
			result: Result{
				StatusCode: 200,
				Size:       1234,
			},
			rules: Rules{
				FilterSize: map[int]bool{
					1234: true,
				},
			},
			want: false,
		},
		{
			name: "matches words",
			result: Result{
				StatusCode: 200,
				Words:      42,
			},
			rules: Rules{
				MatchWords: map[int]bool{
					42: true,
				},
			},
			want: true,
		},
		{
			name: "filters words",
			result: Result{
				StatusCode: 200,
				Words:      42,
			},
			rules: Rules{
				FilterWords: map[int]bool{
					42: true,
				},
			},
			want: false,
		},
		{
			name: "matches lines",
			result: Result{
				StatusCode: 200,
				Lines:      10,
			},
			rules: Rules{
				MatchLines: map[int]bool{
					10: true,
				},
			},
			want: true,
		},
		{
			name: "filters lines",
			result: Result{
				StatusCode: 200,
				Lines:      10,
			},
			rules: Rules{
				FilterLines: map[int]bool{
					10: true,
				},
			},
			want: false,
		},
		{
			name: "matches all conditions",
			result: Result{
				StatusCode: 200,
				Size:       1000,
				Words:      50,
				Lines:      10,
			},
			rules: Rules{
				MatchStatus: map[int]bool{
					200: true,
				},
				MatchSize: map[int]bool{
					1000: true,
				},
				MatchWords: map[int]bool{
					50: true,
				},
				MatchLines: map[int]bool{
					10: true,
				},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Matches(tt.result, tt.rules)

			if got != tt.want {
				t.Fatalf("Matches() = %v, want %v", got, tt.want)
			}
		})
	}
}
