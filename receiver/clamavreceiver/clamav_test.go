package clamavreceiver

import (
	"bufio"
	"errors"
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
			name: "input is present but not include start line",
			text: "Infected files: 1\nTotal errors: 1\nTime: 60.1 sec (1 m )\nEnd Date:   2024:06:26 03:14:37",
			want: Want{-1, -1, -1, nil},
		},
		{
			name: "input is present and within interval",
			text: "----------- SCAN SUMMARY -----------\nInfected files: 1\nTotal errors: 1\nTime: 60.1 sec (1 m )\nEnd Date:   2024:06:26 03:14:37",
			want: Want{1, 1, 60.1, nil},
		},
		{
			name: "input is present and out of interval",
			text: "----------- SCAN SUMMARY -----------\nInfected files: 0\nTotal errors: 0\nTime: 0 sec (0 m )\nEnd Date:   2024:06:26 03:14:34",
			want: Want{-1, -1, -1, nil},
		},
		{
			name: "input is present and same interval",
			text: "----------- SCAN SUMMARY -----------\nInfected files: 0\nTotal errors: 0\nTime: 0 sec (0 m )\nEnd Date:   2024:06:26 03:14:36",
			want: Want{0, 0, 0, nil},
		},
		{
			name: "input is present and within interval",
			text: "Infected files: 2\nTotal errors: 2\nTime: 120.1 sec (2 m )\nEnd Date:   2024:06:26 03:14:37\n----------- SCAN SUMMARY -----------\nInfected files: 1\nTotal errors: 1\nTime: 60.1 sec (1 m )\nEnd Date:   2024:06:26 03:14:37",
			want: Want{1, 1, 60.1, nil},
		},
		{
			name: "input is present first one",
			text: "----------- SCAN SUMMARY -----------\nInfected files: 2\nTotal errors: 2\nTime: 120.1 sec (2 m )\nEnd Date:   2024:06:26 03:14:36\n----------- SCAN SUMMARY -----------\nInfected files: 0\nTotal errors: 0\nTime: 0 sec (0 m )\nEnd Date:   2024:06:26 03:14:37",
			want: Want{2, 2, 120.1, nil},
		},
		{
			name: "infected files invalid",
			text: "----------- SCAN SUMMARY -----------\nInfected files: -\nTotal errors: 1\nTime: 60.1 sec (1 m )\nEnd Date:   2024:06:26 03:14:37",
			want: Want{-1, -1, -1, errors.New("infected files value invalid: strconv.ParseInt: parsing \"-\": invalid syntax")},
		},
		{
			name: "total erros invalid",
			text: "----------- SCAN SUMMARY -----------\nInfected files: 0\nTotal errors: -\nTime: 60.1 sec (1 m )\nEnd Date:   2024:06:26 03:14:37",
			want: Want{-1, -1, -1, errors.New("total errors value invalid: strconv.ParseInt: parsing \"-\": invalid syntax")},
		},
		{
			name: "time invalid",
			text: "----------- SCAN SUMMARY -----------\nInfected files: 0\nTotal errors: 0\nTime: - sec (1 m )\nEnd Date:   2024:06:26 03:14:37",
			want: Want{-1, -1, -1, errors.New("time value invalid: strconv.ParseFloat: parsing \"-\": invalid syntax")},
		},
		{
			name: "end date invalid",
			text: "----------- SCAN SUMMARY -----------\nInfected files: 0\nTotal errors: 0\nTime: 60.1 sec (1 m )\nEnd Date:   -",
			want: Want{-1, -1, -1, errors.New("time value invalid: parsing time \"-\" as \"2006:01:02 15:04:05\": cannot parse \"-\" as \"2006\"")},
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
				t.Errorf("scrape() infectedCount, actual: \n%v\nwant: \n%v", i, tt.want.infectedCount)
			}
			if ttt != tt.want.totalError {
				t.Errorf("scrape() totalError, actual: \n%v\nwant: \n%v", t, tt.want.totalError)
			}
			if e != tt.want.elapsedTime {
				t.Errorf("scrape() elapsedTime, actual: \n%v\nwant: \n%v", e, tt.want.elapsedTime)
			}
			if err != nil || tt.want.err != nil {
				if err.Error() != tt.want.err.Error() {
					t.Errorf("scrape() error, actual: \n%v, want: \n%v", err, tt.want.err)
				}
			}
		})
	}

}
