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
	"io"
	"log"
	"os"

	"github.com/knadh/koanf/providers/posflag"
	"github.com/knadh/koanf/v2"

	flag "github.com/spf13/pflag"
)

const (
	VERSIONstring        = "0.0.2"
	Usage         string = `This is ts, a timestamp tool.

Usage: ts <time string> [<time string>]
-d --diff     Calculate difference between two timestamps (default)
-a --add      Add two timestamps (second parameter must be a duration)
-f --format   For diffs: duration, hour, min, sec, msec
              For timestamps: datetime, rfc3339, date, time, unix, string
              string is a strftime(1) format string. datetime is
              the default
-u --unit    Add unit to the output of timestamp diffs
   --debug    Show debugging output
-v --version  Show program version
-h --help     Show this help screen
-e --examples Show examples or supported inputs`

	Examples string = `Example timestamp inputs:
now                    10:25:30                             Last sunday at 5:30pm                 2 weeks ago
a minute from now      17:25:30                             Next sunday at 22:45                  A week from now
a minute ago           On Friday at noon UTC                Next sunday at 22:45                  A week from today
1 minute ago           On Tuesday at 11am UTC               November 3rd, 1986 at 4:30pm          A month ago
5 minutes ago          On 3 feb 2025 at 5:35:52pm           September 17, 2012 at 10:09am UTC     1 month ago
five minutes ago       3 feb 2025 at 5:35:52pm              September 17, 2012 at 10:09am UTC-8   2 months ago
5 minutes ago          3 days ago at 11:25am                September 17, 2012 at 10:09am UTC+8   12 months ago
2 minutes from now     3 days from now at 14:26             September 17, 2012, 10:11:09          A month from now
two minutes from now   2 weeks ago at 8am                   September 17, 2012, 10:11             One month hence
an hour from now       Today at 10am                        September 17, 2012 10:11              1 month from now
an hour ago            10am today                           September 17 2012 10:11               2 months from now
1 hour ago             Yesterday 10am                       September 17 2012 at 10:11            Last January
6 hours ago            10am yesterday                       Mon Jan  2 15:04:05 2006              Last january
1 hour from now        Yesterday at 10am                    Mon Jan 02 15:04:05 -0700 2006        Next january
noon                   Yesterday at 10:15am                 Mon, 02 Jan 2006 15:04:05 -0700       One year ago
5:35:52pm              Tomorrow 10am                        Mon 02 Jan 2006 15:04:05 -0700        One year from now
10am                   10am tomorrow                        2006-01-02T15:04:05Z                  One year from today
10 am                  Tomorrow at 10am                     1990-12-31T15:59:59-08:00             Two years ago
5pm                    Tomorrow at 10:15am                  One day ago                           2 years ago
10:25am                10:15am tomorrow                     1 day ago                             This year
1:05pm                 Next dec 22nd at 3pm                 3 days ago                            1999AD
10:25:10am             Next December 25th at 7:30am UTC-7   Three days ago                        1999 AD
1:05:10pm              Next December 23rd AT 5:25 PM        1 day from now                        2008CE
10:25                  Last December 23rd AT 5:25 PM        1 week ago                            2008 CE

Example durations for second parameter:
2d1h30m  2 days, one and a half hour
30m      30 minutes`
	ModeDiff int = iota
	ModeAdd
)

type Config struct {
	// TODO: add Timezone parameter
	Showversion bool   `koanf:"version"`
	Debug       bool   `koanf:"debug"`
	Diff        bool   `koanf:"diff"`
	Add         bool   `koanf:"add"`
	Examples    bool   `koanf:"examples"`
	Unit        bool   `koanf:"unit"`
	Format      string `koanf:"format"`
	Args        []string
	Output      io.Writer
	Mode        int
	TZ          string // for unit tests
}

func InitConfig(output io.Writer) (*Config, error) {
	var kloader = koanf.New(".")

	// setup custom usage
	flagset := flag.NewFlagSet("config", flag.ContinueOnError)
	flagset.Usage = func() {
		_, err := fmt.Fprintln(output, Usage)
		if err != nil {
			log.Fatalf("failed to print to output: %s", err)
		}
		os.Exit(0)
	}

	// parse commandline flags
	flagset.BoolP("version", "v", false, "show program version")
	flagset.BoolP("debug", "", false, "enable debug output")
	flagset.BoolP("diff", "d", false, "diff two timestamps")
	flagset.BoolP("add", "a", false, "add two timestamps")
	flagset.BoolP("unit", "u", false, "add unit to diff outputs")
	flagset.BoolP("examples", "e", false, "show examples of supported inputs")
	flagset.StringP("format", "f", "", "format to print timestamps or diffs")

	if err := flagset.Parse(os.Args[1:]); err != nil {
		return nil, fmt.Errorf("failed to parse program arguments: %w", err)
	}

	// command line setup
	if err := kloader.Load(posflag.Provider(flagset, ".", kloader), nil); err != nil {
		return nil, fmt.Errorf("error loading flags: %w", err)
	}

	// fetch values
	conf := &Config{Output: output}
	if err := kloader.Unmarshal("", &conf); err != nil {
		return nil, fmt.Errorf("error unmarshalling: %w", err)
	}

	// want examples?
	if conf.Examples {
		_, err := fmt.Fprintln(output, Examples)
		if err != nil {
			Die("failed write to output file handle", err)
		}

		os.Exit(0)
	}

	// args are timestamps
	if len(flagset.Args()) == 0 {
		return nil, errors.New("no timestamp argument[s] specified.\n" + Usage)
	}

	conf.Args = flagset.Args()

	if conf.Add {
		conf.Mode = ModeAdd
	} else {
		conf.Mode = ModeDiff
	}

	return conf, nil
}
