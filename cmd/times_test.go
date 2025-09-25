/*
Copyright © 2025 Thomas von Dein

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// reference time for our tests
var now = time.Date(2025, 9, 25, 12, 30, 00, 0, time.UTC)

// return fixed date at given params
func dateAtTime(dateFrom time.Time, hour int, min int, sec int) time.Time {
	t := dateFrom
	return time.Date(t.Year(), t.Month(), t.Day(), hour, min, sec, 0, t.Location())
}

func setUTC() {
	// make sure to have a consistent environment
	if err := os.Setenv("TZ", "UTC"); err != nil {
		log.Fatal(err)
	}

	loc, _ := time.LoadLocation("UTC")
	time.Local = loc
}

func TestParseTimestamps(t *testing.T) {
	setUTC()

	var datetimes = []struct {
		input string
		want  time.Time
	}{
		// some timestamps from ijt/go-anytime/anytime_test.go
		{`a minute from now`, now.Add(time.Minute)},
		{`5 minutes ago`, now.Add(-5 * time.Minute)},
		{`an hour from now`, now.Add(time.Hour)},
		{`Yesterday 10am`, dateAtTime(now.AddDate(0, 0, -1), 10, 0, 0)},
		{`Mon Jan  2 15:04:05 2006`, time.Date(2006, 1, 2, 15, 4, 5, 0, now.Location())},
		{`1 day from now`, now.Add(24 * time.Hour)},
		{`One year ago`, now.AddDate(-1, 0, 0)},
		{`03:15`, dateAtTime(now, 3, 15, 0)},
		{`Wed Sep 25 12:30:00 PM UTC 2025`, now},
		{`Wed Sep 25 00:30:00 PM CEST 2025`, now},
		{`Wed Sep 25 2025 13:30:00 GMT+0100 (GMT Daylight Time)`, now},
	}

	for _, tt := range datetimes {
		testname := fmt.Sprintf("parsetimestamp-%s", strings.ReplaceAll(tt.input, " ", "-"))
		t.Run(testname, func(t *testing.T) {
			var writer bytes.Buffer
			tp := NewTP(&Config{Args: []string{tt.input}, Output: &writer, tz: "UTC"}, now)

			// writer.String()
			ts, err := tp.Parse(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, ts)

			err = tp.ProcessTimestamps()
			assert.NoError(t, err)
			assert.EqualValues(t, tt.want.Format(DefaultFormat)+"\n", writer.String())
		})
	}
}

func preDiff(tsA, tsB time.Time) time.Duration {
	var diff time.Duration

	// avoid negative results
	if tsA.Unix() > tsB.Unix() {
		diff = tsA.Sub(tsB)
	} else {
		diff = tsB.Sub(tsA)
	}

	return diff
}

func TestDiffTimestamps(t *testing.T) {
	setUTC()

	var datetimes = []struct {
		A    string
		B    string
		want time.Duration
	}{
		{`now`, `11:30`, preDiff(now, dateAtTime(now, 11, 30, 00))},
		{`11:30`, `now`, preDiff(now, dateAtTime(now, 11, 30, 00))},
	}

	for _, tt := range datetimes {
		testname := fmt.Sprintf("diff-%s-%s", strings.ReplaceAll(tt.A, " ", "-"), strings.ReplaceAll(tt.B, " ", "-"))
		t.Run(testname, func(t *testing.T) {
			var writer bytes.Buffer
			tp := NewTP(&Config{Args: []string{tt.A, tt.B}, Output: &writer, Mode: ModeDiff}, now)

			err := tp.ProcessTimestamps()
			assert.NoError(t, err)
			assert.EqualValues(t, tt.want.String()+"\n", writer.String())
		})
	}
}

func TestAddTimestamps(t *testing.T) {
	setUTC()

	var datetimes = []struct {
		A    string
		B    string
		want time.Time
	}{
		{`now`, `01:30`, dateAtTime(now, 14, 00, 00)},
		{`now`, `2h`, dateAtTime(now, 14, 30, 00)},
		{`now`, `12d4h`, dateAtTime(now.Add(time.Hour*24*12), 16, 30, 00)},
		{`now`, `45m`, dateAtTime(now, 13, 15, 00)},
		{`now`, `1d10s`, dateAtTime(now.Add(time.Hour*24*1), 12, 30, 10)},
	}

	for _, tt := range datetimes {
		testname := fmt.Sprintf("diff-%s-%s", strings.ReplaceAll(tt.A, " ", "-"), strings.ReplaceAll(tt.B, " ", "-"))
		t.Run(testname, func(t *testing.T) {
			var writer bytes.Buffer
			tp := NewTP(&Config{Args: []string{tt.A, tt.B}, Output: &writer, Mode: ModeAdd}, now)

			err := tp.ProcessTimestamps()
			assert.NoError(t, err)
			assert.EqualValues(t, tt.want.Format(DefaultFormat)+"\n", writer.String())
		})
	}
}
