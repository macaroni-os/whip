/*

Copyright (C) 2021  Daniele Rondina <geaaru@sabayonlinux.org>

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
package executor

import (
	"fmt"
	"io"
	"os"
)

type WhipWriter struct {
	Type string
}

func NewWhipWriter(t string) *WhipWriter {
	return &WhipWriter{
		Type: t,
	}
}

func (e *WhipWriter) Write(p []byte) (int, error) {
	var wr io.Writer

	if e.Type == "stdout" {
		wr = os.Stdout
	} else {
		wr = os.Stderr
	}

	fmt.Fprint(wr, string(p))

	/*
		logger.Msg("info", false, false,
			logger.Aurora.Bold(
				logger.Aurora.BrightCyan(string(p)),
			),
		)
	*/
	return len(p), nil
}

func (e *WhipWriter) Close() error {
	return nil
}
