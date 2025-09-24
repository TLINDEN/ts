# ts

generic cli timestamp parser and calculator tool

## Usage

```default
This is ts, a timestamp tool.

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

## Examples

```default
# diff between to day and yesterday 10 am
ts today "10am yesterday"
14h0m0s

# show timestamp from a couple days ago
ts "3 days ago"
2025-09-21 13:51:55.054754744 +0200 CEST

# show timestamp of one hour and 45 minutes before (-d is the defaul)
ts -d now 1h45m
2025-09-24 12:07:45.072300157 +0200 CEST m=-6299.999710536

# 10 hours from now
ts now 10h
2025-09-24 03:53:36.7095512 +0200 CEST m=-35999.999720767
```

To see a comprehensive list of supported inputs, call `ts -e`.

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
