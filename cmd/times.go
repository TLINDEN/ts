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
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/araddon/dateparse"
	"github.com/ijt/go-anytime"
	modnow "github.com/jinzhu/now"
)

type TimestampProccessor struct {
	Config
	Reference time.Time
}

func NewTP(conf *Config, ref ...time.Time) *TimestampProccessor {
	// we add some pre-defined formats to modnow
	formats := []string{
		time.UnixDate, time.RubyDate,
		time.RFC1123, time.RFC1123Z, time.RFC3339, time.RFC3339Nano,
		time.RFC822, time.RFC822Z, time.RFC850,
		"Mon Jan 02 15:04:05 PM MST 2006", // linux date
	}

	modnow.TimeFormats = append(modnow.TimeFormats, formats...)

	tp := &TimestampProccessor{Config: *conf, Reference: time.Now()}

	if len(ref) == 1 {
		tp.Reference = ref[0]
	}

	return tp
}

func (tp *TimestampProccessor) ProcessTimestamps() error {
	switch len(tp.Args) {
	case 1:
		return tp.SingleTimestamp(tp.Args[0])
	case 2:
		return tp.DualTimestamps(tp.Args[0], tp.Args[1])
	}

	return nil
}

// a post processor for ParseTimestamp() to apply custom time zone, if any
func (tp *TimestampProccessor) Parse(timestamp string) (time.Time, error) {
	ts, err := tp.ParseTimestamp(timestamp)

	if err != nil {
		return ts, err
	}

	if tp.tz != "" {
		// apply custom timezone
		zone, _ := time.LoadLocation(tp.tz)
		ts = ts.In(zone)
	}

	return ts, nil
}

// Parse uses 3 different timestamp parser modules to provide maximum flexibility
func (tp *TimestampProccessor) ParseTimestamp(timestamp string) (time.Time, error) {
	ts, err := anytime.Parse(timestamp, tp.Reference)
	if err == nil {
		return ts, nil
	}

	// anytime failed, try module modnow
	ts, err = modnow.Parse(timestamp)
	if err == nil {
		return ts, nil
	}

	// modnow failed, try module dateparse
	return dateparse.ParseAny(timestamp)
}

func (tp *TimestampProccessor) SingleTimestamp(timestamp string) error {
	ts, err := tp.Parse(timestamp)
	if err != nil {
		return err
	}

	tp.Print(TPdatetime{TimestampProccessor: *tp, Data: ts})
	//tp.Print(ts)

	return nil
}

func (tp *TimestampProccessor) DualTimestamps(timestampA, timestampB string) error {
	tsA, err := tp.Parse(timestampA)
	if err != nil {
		return err
	}

	durB, err := duration2int(timestampB)
	if err == nil {
		// calculate with a duration
		tp.CalcDuration(tsA, durB)
		return nil
	}

	tsB, err := tp.Parse(timestampB)
	if err != nil {
		return err
	}

	tp.CalcDiff(tsA, tsB)

	return nil
}

func (tp *TimestampProccessor) CalcDiff(tsA time.Time, tsB time.Time) {
	switch tp.Mode {
	case ModeDiff:
		var diff time.Duration

		// avoid negative results
		if tsA.Unix() > tsB.Unix() {
			diff = tsA.Sub(tsB)
		} else {
			diff = tsB.Sub(tsA)
		}

		tp.Print(TPduration{TimestampProccessor: *tp, Data: diff})

	case ModeAdd:
		seconds := (tsB.Hour() * 3600) + (tsB.Minute() * 60) + tsB.Second()
		sum := tsA.Add(time.Duration(seconds) * time.Second)

		tp.Print(TPdatetime{TimestampProccessor: *tp, Data: sum})
	}
}

func (tp *TimestampProccessor) CalcDuration(tsA time.Time, durB time.Duration) {
	var datetime time.Time

	switch tp.Mode {
	case ModeDiff:
		datetime = tsA.Add(-durB)
	case ModeAdd:
		datetime = tsA.Add(durB)
	}

	tp.Print(TPdatetime{TimestampProccessor: *tp, Data: datetime})
}

func (tp *TimestampProccessor) Print(ts TimestampWriter) {
	_, err := fmt.Fprintln(tp.Output, ts.String())
	if err != nil {
		Die("failed to print to given output handle", err)
	}
}

/*
We could use time.ParseDuration(), but this doesn't support days.

We  could also  use github.com/xhit/go-str2duration/v2,  which does
the job,  but it's  just another dependency,  just for  this little
gem. And  we don't need a  time.Time value.

Convert a  duration into  seconds (int).

	Valid  time units  are "s", "m", "h" and "d".

	Valid inputs: 2h5m (2 hours and 5 min), 10d12h (10 and a half days)
*/
func duration2int(duration string) (time.Duration, error) {
	re := regexp.MustCompile(`(\d+)([dhms])`)
	seconds := 0
	found := false

	for _, match := range re.FindAllStringSubmatch(duration, -1) {
		if len(match) == 3 {
			found = true
			v, _ := strconv.Atoi(match[1])
			switch match[2][0] {
			case 'd':
				seconds += v * 86400
			case 'h':
				seconds += v * 3600
			case 'm':
				seconds += v * 60
			case 's':
				seconds += v
			}
		}
	}

	if !found {
		return 0, errors.New("failed to parse duration")
	}

	return time.Duration(seconds) * time.Second, nil
}
