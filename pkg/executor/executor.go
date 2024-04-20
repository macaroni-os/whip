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
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"

	log "github.com/macaroni-os/whip/pkg/logger"
	specs "github.com/macaroni-os/whip/pkg/specs"
)

type Executor struct {
	Config *specs.Config
}

func NewExecutor(config *specs.Config) *Executor {
	return &Executor{
		Config: config,
	}
}

func (e *Executor) RunCommandWithOutput(command string, envs map[string]string,
	outBuffer, errBuffer io.WriteCloser, entryPoint []string) (int, error) {
	ans := 1

	if len(entryPoint) == 0 {
		entryPoint = []string{"/bin/bash", "-c"}
	}

	if outBuffer == nil {
		return 1, errors.New("Invalid outBuffer")
	}
	if errBuffer == nil {
		return 1, errors.New("Invalid errBuffer")
	}

	cmds := append(entryPoint, command)

	cmd := exec.Command(cmds[0], cmds[1:]...)

	logger := log.GetDefaultLogger()

	logger.DebugC(
		fmt.Sprintf(":rocket: - %s - [%s]",
			entryPoint, command),
	)

	// Convert envs to array list
	elist := os.Environ()
	for k, v := range envs {
		elist = append(elist, k+"="+v)
	}

	cmd.Stdout = outBuffer
	cmd.Stderr = errBuffer
	cmd.Env = elist

	err := cmd.Start()
	if err != nil {
		logger.Error("Error on start command: " + err.Error())
		return 1, err
	}

	err = cmd.Wait()
	if err != nil {
		logger.Error("Error on waiting command: " + err.Error())
		return 1, err
	}

	ans = cmd.ProcessState.ExitCode()

	logger.DebugC(
		fmt.Sprintf(":airplane_arriving: - Exiting [%d]", ans),
	)

	return ans, nil
}
