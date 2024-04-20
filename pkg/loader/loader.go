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
package loader

import (
	"errors"
	"fmt"
	"strings"

	exec "github.com/macaroni-os/whip/pkg/executor"
	log "github.com/macaroni-os/whip/pkg/logger"
	specs "github.com/macaroni-os/whip/pkg/specs"
)

type Loader interface {
	Load() (map[string]*specs.SpecFile, error)
}

type WhipHolder struct {
	Config *specs.Config
	Logger *log.Logger
	Rules  map[string]*specs.SpecFile
	Envs   map[string]string
}

func NewWhipHolder(config *specs.Config) *WhipHolder {
	ans := &WhipHolder{
		Config: config,
		Logger: log.NewLogger(config),
		Rules:  make(map[string]*specs.SpecFile, 0),
		Envs:   make(map[string]string, 0),
	}

	// Initialize logging
	if config.GetLogging().EnableLogFile && config.GetLogging().Path != "" {
		err := ans.Logger.InitLogger2File()
		if err != nil {
			ans.Logger.Fatal("Error on initialize logfile")
		}
	}
	ans.Logger.SetAsDefault()
	return ans
}

func (w *WhipHolder) GetConfig() *specs.Config              { return w.Config }
func (w *WhipHolder) GetLogger() *log.Logger                { return w.Logger }
func (w *WhipHolder) GetRules() *map[string]*specs.SpecFile { return &w.Rules }

func (w *WhipHolder) GetRule(s string) (*specs.SpecFile, error) {
	ans, ok := w.Rules[s]
	if ok {
		return ans, nil
	}

	return nil, errors.New(fmt.Sprintf("No rule found with name %s", s))
}

func (w *WhipHolder) RunHook(h string) error {
	dot := strings.Index(h, ".")

	if dot < 0 {
		return errors.New("Filename not specified on hook " + h)
	}

	filename := h[0:dot]
	hname := h[dot+1:]

	specFile, err := w.GetRule(filename)
	if err != nil {
		return err
	}

	hook, err := specFile.GetHook(hname)
	if err != nil {
		return err
	}

	if len(hook.GetActions()) == 0 {
		return errors.New(fmt.Sprintf(
			"Hook %s without actions!", h,
		))
	}

	w.Logger.DebugC(fmt.Sprintf(
		"Running hook %s of the file %s.", hname, filename))

	wso := exec.NewWhipWriter("stdout")
	wse := exec.NewWhipWriter("stderr")
	for _, command := range hook.GetActions() {
		executor := exec.NewExecutor(w.Config)
		res, err := executor.RunCommandWithOutput(
			command,
			w.Envs,
			wso, wse,
			specFile.GetEntrypoint(),
		)

		if err != nil {
			return errors.New(
				fmt.Sprintf("[%s] Error - %s", h, err.Error()),
			)
		} else if res != 0 {
			return errors.New(
				fmt.Sprintf("[%s] Exiting with %d", h, res),
			)
		}
	}

	w.Logger.InfoC(fmt.Sprintf(
		"[%s] Completed correctly.", h))

	return nil
}

func (w *WhipHolder) LoadSpecs() error {
	var loader Loader

	switch w.Config.Loader {
	case "dir":
		loader = NewLoaderDir(w.Config)
	default:
		return errors.New("Invalid loader defined")
	}

	m, err := loader.Load()
	if err != nil {
		return err
	}

	w.Rules = m

	w.Logger.DebugC(
		":factory: Loaded", len(w.Rules), "files.",
	)

	return nil
}
