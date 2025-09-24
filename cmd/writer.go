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
	"time"
)

type TimestampWriter interface {
	String() string
}

type TPduration struct {
	TimestampProccessor
	Data time.Duration
}

type TPdatetime struct {
	TimestampProccessor
	Data time.Time
}

func (duration TPduration) String() string {
	var unit string

	if duration.Unit {
		switch duration.Format {
		case "d", "day", "days":
			unit = " days"
		case "h", "hour", "hours":
			unit = " hours"
		case "m", "min", "mins", "minutes":
			unit = " minutes"
		case "s", "sec", "secs", "seconds":
			unit = " seconds"
		case "ms", "msec", "msecs", "milliseconds":
			unit = " milliseconds"
		}
	}

	// duration, days, hour, min, sec, msec
	switch duration.Format {
	case "d", "day", "days":
		return fmt.Sprintf("%.02f%s", duration.Data.Hours()/24, unit)
	case "h", "hour", "hours":
		return fmt.Sprintf("%.02f%s", duration.Data.Hours(), unit)
	case "m", "min", "mins", "minutes":
		return fmt.Sprintf("%.02f%s", duration.Data.Minutes(), unit)
	case "s", "sec", "secs", "seconds":
		return fmt.Sprintf("%.02f%s", duration.Data.Seconds(), unit)
	case "ms", "msec", "msecs", "milliseconds":
		return fmt.Sprintf("%d%s", duration.Data.Milliseconds(), unit)
	case "dur", "duration":
		fallthrough
	default:
		return duration.Data.String()
	}
}

func (datetime TPdatetime) String() string {
	// datetime(default), date, time, unix, string
	switch datetime.Format {
	case "rfc3339":
		return datetime.Data.Format(time.RFC3339)
	case "date":
		return datetime.Data.Format("2006-01-02")
	case "time":
		return datetime.Data.Format("03:04:05")
	case "unix":
		return fmt.Sprintf("%d", datetime.Data.Unix())
	case "datetime":
		fallthrough
	case "":
		return datetime.Data.String()
	default:
		return datetime.Data.Format(datetime.Format)
	}
}
