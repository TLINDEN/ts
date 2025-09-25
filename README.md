# ts

generic cli timestamp parser and calculator tool

## Introduction

This little utility is a commandline frontent to the amazing datetime
parser module [anytime](https://github.com/ijt/go-anytime). It uses
two other modules as fallback if anytime might fail:
[now](https://github.com/jinzhu/now) and
[dateparse](github.com/araddon/dateparse).

You can use it to print timestamps from plain english phrases like
`next December 23rd AT 5:25 PM` or `two minutes from now`. In addition
you can calculate the difference between two timestamps and you can
add a duration to a timestamp.


## Example Usage

In these examples the current time is always **2025-09-17T07:30:00+01:00**.

Show current date and time (same as `date`):
```default
% ts now
Wed Sep 17 07:30:00 +0100 2025
```

show timestamp for minus 1 hour
```default
% ts "1 hour ago"
Wed Sep 17 06:30:00 +0100 2025
```

... or from a couple days ago:
```default
% ts "4 days ago"
Sat Sep 13 07:30:00 +0100 2025
```
There are much more ways to get timestamps, see `ts -e`.

We  can  also add  times  to  timestamps, here  we  want  to know  the
timestamp from now plus 10 days and 4 hours in the future:
```default
% ts -a now 10d4h
Sat Sep 27 11:30:00 +0100 2025
```

It doesn't make a difference where you position the `-a` parameter:
```default
% ts now -a 10d4h
Sat Sep 27 11:30:00 +0100 2025
```
Of course you can also calculate the difference between two
dates. Here we have two timestamps (maybe we took them from a log
file) and want to know the dime elapsed between them:

```default
% ts 2025-09-17T07:30:00+01:00 2025-09-15T12:45:00+01:00
42h45m0s
```
As you can see, if you do not provide a parameter, the default is to
calculate the difference between the two args. To explicitly calculate
the difference, use the `-d` parameter.

You can of course use english phrases for time differences as well:
```default
% ts "today 9 am" 2025-09-15T12:45:00+01:00
44h15m0s
```

Lets talk a little bit about formatting. You may have already
recognized, that `ts` prints either whole timestamps or
durations. Both output types can be modified with the `-f`
parameter. There are predefined formats for timestamps:

```default
% ts now 
Wed Sep 17 07:30:00 +0100 2025
% ts now -f rfc3339
2025-09-17T07:30:00+01:00
% ts now -f date
2025-09-17
% ts now -f unix
1758090600
```

But you can also specify your own, you have to follow the [golang
rules for timestamp formats](https://pkg.go.dev/time#Layout),
basically:

* Year: "2006" "06"
* Month: "Jan" "January" "01" "1"
* Day of the week: "Mon" "Monday"
* Day of the month: "2" "_2" "02"
* Day of the year: "__2" "002"
* Hour: "15" "3" "03" (PM or AM)
* Minute: "4" "04"
* Second: "5" "05"
* AM/PM mark: "PM"

for example:
```default
% ts now -f "Mon, 02.January 2006"
Wed, 17.September 2025
```

Ok I admit look is kinda weird, complaints go the the golang dev team
:).

Duration formatting is also customizable. By default a duration looks
like we have seen above: `44h15m0s`. But sometimes we want to know the
number of hours or minutes. Easy:

```default
% ts now 2025-09-15T12:45:00+01:00 -f hours
42.75
% ts now 2025-09-15T12:45:00+01:00 -f minutes
2565.00
```

You may also add the `-u` parameter to have the unit shown as well:

```default
% ts now 2025-09-15T12:45:00+01:00 -f hours -u
42.75 hours
% ts now 2025-09-15T12:45:00+01:00 -f minutes -u
2565.00 minutes
```

## Commandline parameters

Here is the list of all supported parameters:

```default
Usage: ts <time string> [<time string>]
-d --diff     Calculate difference between two timestamps (default).
-a --add      Add two timestamps (second parameter must be a time).
-f --format   For diffs: duration, hour, min, sec, msec.
              For timestamps: datetime, rfc3339, date, time, unix, string.
              string is a strftime(1) format string. datetime is
              the default.
-u --unit    Add unit to the output of timestamp diffs.
   --debug    Show debugging output.
-v --version  Show program version.
-h --help     Show this help screen.
-e --examples Show examples or supported inputs.
```


## Installation

The tool does not have any dependencies.  Just download the binary for
your platform from the releases page and you're good to go.

### Installation using a pre-compiled binary

Go to the [latest release page](https://github.com/TLINDEN/ts/releases/latest)
and look for your OS and platform. There are two options to install the binary:

Directly     download     the     binary    for     your     platform,
e.g. `ts-linux-amd64-0.0.2`, rename it to `ts` (or whatever
you like more!)  and put it into  your bin dir (e.g. `$HOME/bin` or as
root to `/usr/local/bin`).

Be sure  to verify  the signature  of the binary  file. For  this also
download the matching `ts-linux-amd64-0.0.2.sha256` file and:

```shell
cat ts-linux-amd64-0.0.2.sha25 && sha256sum ts-linux-amd64-0.0.2
```
You should see the same SHA256 hash.

You  may  also download  a  binary  tarball  for your  platform,  e.g.
`ts-linux-amd64-0.0.2.tar.gz`,  unpack and  install it.  GNU Make  is
required for this:
   
```shell
tar xvfz ts-linux-amd64-0.0.2.tar.gz
cd ts-linux-amd64-0.0.2
sudo make install
```

### Installation from source

You will need the Golang toolchain  in order to build from source. GNU
Make will also help but is not strictly neccessary.

If you want to compile the tool yourself, use `git clone` to clone the
repository.   Then   execute   `go    mod   tidy`   to   install   all
dependencies. Then  just enter `go  build` or -  if you have  GNU Make
installed - `make`.

To install after building either copy the binary or execute `sudo make
install`. 

# Report bugs

[Please open an issue](https://github.com/TLINDEN/ts/issues). Thanks!

# License

This work is licensed under the terms of the General Public Licens
version 3.

# Author

Copyleft (c) 2025 Thomas von Dein
