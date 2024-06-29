package clamavreceiver

import (
	"bufio"
	"strings"
	"testing"
	"time"

	"go.opentelemetry.io/collector/pdata/pcommon"
)

func TestScrape(t *testing.T) {
	type Want struct {
		infectedCount int64
		totalError    int64
		elapsedTime   float64
		err           error
	}
	tests := []struct {
		name string
		text string
		want Want
	}{
		{
			name: "input is nil",
			text: "",
			want: Want{-1, -1, -1, nil},
		},
		{
			name: "input is present",
			text: "Infected files: 0\nTotal errors: 1\nTime: 60.1 sec (1 m )\nEnd Date:   2024:06:26 03:14:44",
			want: Want{0, 1, 60.1, nil},
		},
		{
			name: "input is present but not within interval",
			text: "Infected files: 0\nTotal errors: 1\nTime: 60.1 sec (1 m )\nEnd Date:   2024:06:26 03:03:59",
			want: Want{-1, -1, -1, nil},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Helper()
			reader := strings.NewReader(tt.text)
			scanner := bufio.NewScanner(reader)
			d, _ := time.ParseDuration("10s")
			now := time.Date(2024, 6, 26, 3, 14, 45, 11, &time.Location{})
			i, ttt, e, err := scrape(scanner, d, pcommon.NewTimestampFromTime(now))
			if i != tt.want.infectedCount {
				t.Errorf("scrape() infectedCount, actual: %v, want: %v", i, tt.want.infectedCount)
			}
			if ttt != tt.want.totalError {
				t.Errorf("scrape() totalError, actual: %v, want: %v", t, tt.want.totalError)
			}
			if e != tt.want.elapsedTime {
				t.Errorf("scrape() elapsedTime, actual: %v, want: %v", e, tt.want.elapsedTime)
			}
			if err != nil || tt.want.err != nil {
				if err.Error() != tt.want.err.Error() {
					t.Errorf("scrape() error, actual: %v, want %v", err, tt.want.err)
				}
			}
		})
	}

}
