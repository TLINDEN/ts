package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/ijt/go-anytime"
	"github.com/itlightning/dateparse"
	modnow "github.com/jinzhu/now"
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

	modnow.TimeFormats = append(modnow.TimeFormats, formats...)

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

	// anytime failed, try module modnow
	ts, err = modnow.Parse(timestamp)
	if err == nil {
		return ts, nil
	}

	// modnow failed, try module dateparse
	return dateparse.ParseAny(timestamp)
}

func (tp *TimestampProccessor) Calc(timestampA, timestampB string) error {
	tsA, err := tp.Parse(timestampA)
	if err != nil {
		return err
	}

	tsB, err := tp.Parse(timestampB)
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
		tp.Print(TPduration{TimestampProccessor: *tp, Data: diff})

	case ModeAdd:
		seconds := (tsB.Hour() * 3600) + (tsB.Minute() * 60) + tsB.Second()
		sum := tsA.Add(time.Duration(seconds) * time.Second)
		tp.Print(TPdatetime{TimestampProccessor: *tp, Data: sum})
	}

	return nil
}

func (tp *TimestampProccessor) Print(ts TimestampWriter) {
	_, err := fmt.Fprintln(tp.Output, ts.String())
	if err != nil {
		log.Fatalf("failed to print to given output handle: %s", err)
	}
}
