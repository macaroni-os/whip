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
	"io/ioutil"
	"path"
	"path/filepath"
	"regexp"

	log "github.com/macaroni-os/whip/pkg/logger"
	specs "github.com/macaroni-os/whip/pkg/specs"
)

// LoaderDir permits to laod whip specs file
// from a list of directories defined.
type LoaderDir struct {
	Config *specs.Config
}

func NewLoaderDir(c *specs.Config) *LoaderDir {
	return &LoaderDir{
		Config: c,
	}
}

func (l *LoaderDir) Load() (map[string]*specs.SpecFile, error) {
	var regexConfs = regexp.MustCompile(`.yml$|.yaml$`)

	logger := log.GetDefaultLogger()

	ans := make(map[string]*specs.SpecFile, 0)
	if len(l.Config.SpecDirs) == 0 {
		return ans, errors.New("No specs dirs defined!")
	}

	for _, sdir := range l.Config.SpecDirs {
		logger.Debug("Checking directory", sdir, "...")

		files, err := ioutil.ReadDir(sdir)
		if err != nil {
			logger.Debug("Skip dir", sdir, ":", err.Error())
			continue
		}

		for _, file := range files {
			if file.IsDir() {
				continue
			}

			if !regexConfs.MatchString(file.Name()) {
				logger.Debug("File", file.Name(), "skipped.")
				continue
			}

			sfile := path.Join(sdir, file.Name())
			content, err := ioutil.ReadFile(sfile)
			if err != nil {
				logger.Debug("On read file", file.Name(), ":", err.Error())
				logger.Debug("File", file.Name(), "skipped.")
				continue
			}

			sf, err := specs.NewSpecFileFromYaml(content, sfile)
			if err != nil {
				return ans, err
			}

			ext := filepath.Ext(file.Name())
			key := file.Name()
			if len(ext) > 0 {
				key = key[0 : len(key)-len(ext)]
			}

			ans[key] = sf
		}

	}

	return ans, nil
}
