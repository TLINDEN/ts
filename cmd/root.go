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
	"io"
	"os"
)

func Die(format string, err error) int {
	fmt.Fprintf(os.Stderr, format+": %s\n", err)

	return 1
}

func Main(output io.Writer) int {
	conf, err := InitConfig(output)
	if err != nil {
		fmt.Println(1)
		return Die("failed to initialize", err)
	}

	tp := NewTP(conf)

	if err := tp.ProcessTimestamps(); err != nil {
		return Die("failed to process timestamp[s]", err)
	}

	return 0
}
