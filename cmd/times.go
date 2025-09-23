package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/ijt/go-anytime"
	"github.com/itlightning/dateparse"
	"github.com/jinzhu/now"
)

type TimestampProccessor struct {
	Config
}

func NewTP(conf *Config) *TimestampProccessor {
	formats := []string{
		time.UnixDate, time.RubyDate,
		time.RFC1123, time.RFC1123Z, time.RFC3339, time.RFC3339Nano,
		time.RFC822, time.RFC822Z, time.RFC850,
		"Mon Jan 02 15:04:05 PM MST 2006", // linux date
		"Mo. 02 Jan. 2006 15:04:05 MST",   // freebsd date (fails, see golang/go/issues/75576)
	}

	now.TimeFormats = append(now.TimeFormats, formats...)

	return &TimestampProccessor{Config: *conf}
}

func (tp *TimestampProccessor) ProcessTimestamps() error {
	switch len(tp.Args) {
	case 1:
		return tp.SingleTimestamp(tp.Args[0])
	case 2:
		return tp.Calc(tp.Args[0], tp.Args[1])
	}

	return nil
}

func (tp *TimestampProccessor) SingleTimestamp(timestamp string) error {
	ts, err := tp.Parse(timestamp)
	if err != nil {
		return err
	}

	tp.Print(ts)

	return nil
}

// Parse uses 3 different timestamp parser modules to provide the maximum flexibility
func (tp *TimestampProccessor) Parse(timestamp string) (time.Time, error) {
	reference := time.Now()
	ts, err := anytime.Parse(timestamp, reference)
	if err == nil {
		return ts, nil
	}

	// anytime failed, try module now
	ts, err = now.Parse(timestamp)
	if err == nil {
		return ts, nil
	}

	// now failed, try module dateparse
	ts, err = dateparse.ParseAny(timestamp)

	return ts, nil
}

func (tp *TimestampProccessor) Calc(timestampA, timestampB string) error {
	now := time.Now()
	tsA, err := anytime.Parse(timestampA, now)
	if err != nil {
		return err
	}

	tsB, err := anytime.Parse(timestampB, now)
	if err != nil {
		return err
	}

	switch tp.Mode {
	case ModeDiff:
		var diff time.Duration
		if tsA.Unix() > tsB.Unix() {
			diff = tsA.Sub(tsB)
		} else {
			diff = tsB.Sub(tsA)
		}
		tp.Print(diff)
	case ModeAdd:
		seconds := (tsB.Hour() * 3600) + (tsB.Minute() * 60) + tsB.Second()
		tp.Print(tsA.Add(time.Duration(seconds) * time.Second))
	}

	return nil
}

func (tp *TimestampProccessor) Print(msg any) {
	var repr string

	switch msg := msg.(type) {
	case string:
		repr = msg
	case time.Time:
		repr = tp.StringTime(msg)
	case time.Duration:
		repr = tp.StringDuration(msg)
	}

	_, err := fmt.Fprintln(tp.Output, repr)
	if err != nil {
		log.Fatalf("failed to print to given output handle: %s", err)
	}
}

func (tp *TimestampProccessor) StringDuration(msg time.Duration) string {
	var unit string

	if tp.Unit {
		switch tp.Format {
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
	switch tp.Format {
	case "d", "day", "days":
		return fmt.Sprintf("%.02f%s", msg.Hours()/24+(msg.Minutes()/60), unit)
	case "h", "hour", "hours":
		return fmt.Sprintf("%.02f%s", msg.Hours(), unit)
	case "m", "min", "mins", "minutes":
		return fmt.Sprintf("%.02f%s", msg.Minutes(), unit)
	case "s", "sec", "secs", "seconds":
		return fmt.Sprintf("%.02f%s", msg.Seconds(), unit)
	case "ms", "msec", "msecs", "milliseconds":
		return fmt.Sprintf("%d%s", msg.Milliseconds(), unit)
	case "dur", "duration":
		fallthrough
	default:
		return msg.String()
	}
}

func (tp *TimestampProccessor) StringTime(msg time.Time) string {
	// datetime(default), date, time, unix, string
	switch tp.Format {
	case "rfc3339":
		return msg.Format(time.RFC3339)
	case "date":
		return msg.Format("2006-01-02")
	case "time":
		return msg.Format("03:04:05")
	case "unix":
		return fmt.Sprintf("%d", msg.Unix())
	case "datetime":
		fallthrough
	case "":
		return msg.String()
	default:
		return msg.Format(tp.Format)
	}
}
