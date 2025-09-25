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
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDuration(t *testing.T) {
	var tests = []struct {
		format   string
		duration time.Duration
		want     string
	}{
		{"day", time.Hour * 24, "1.00 days"},
		{"dur", time.Hour * 24, "24h0m0s"},
		{"hour", time.Hour * 24, "24.00 hours"},
		{"min", time.Minute * 20, "20.00 minutes"},
		{"min", time.Minute*20 + time.Second*30, "20.50 minutes"},
		{"sec", time.Second * 30, "30.00 seconds"},
		{"ms", time.Second * 30, "30000 milliseconds"},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("formatduration-%s", tt.format)
		t.Run(testname, func(t *testing.T) {
			tpdur := TPduration{Data: tt.duration}
			tpdur.Unit = true
			tpdur.Format = tt.format

			out := tpdur.String()
			assert.Equal(t, tt.want, out)
		})
	}
}

func TestDatetime(t *testing.T) {
	var now = time.Date(2025, 9, 25, 12, 30, 00, 0, time.UTC)

	var tests = []struct {
		format   string
		datetime time.Time
		want     string
	}{
		{"rfc3339", now, "2025-09-25T12:30:00Z"},
		{"date", now, "2025-09-25"},
		{"time", now, "12:30:00"},
		{"unix", now, "1758803400"},
		{"datetime", now, "Thu Sep 25 12:30:00 UTC 2025"},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("formatdatetime-%s", tt.format)
		t.Run(testname, func(t *testing.T) {
			tpdat := TPdatetime{Data: tt.datetime}
			tpdat.Format = tt.format

			out := tpdat.String()
			assert.Equal(t, tt.want, out)
		})
	}
}
