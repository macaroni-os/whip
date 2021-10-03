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
package specs

import (
	"errors"
	"fmt"

	"gopkg.in/yaml.v2"
)

func (f *SpecFile) GetHook(s string) (Hook, error) {
	ans, ok := f.Hooks[s]
	if ok {
		return ans, nil
	}

	return Hook{}, errors.New(fmt.Sprintf("No hook found with name %s", s))
}

func (f *SpecFile) GetEntrypoint() []string {
	if len(f.Entrypoint) > 0 {
		return f.Entrypoint
	}
	// Using /bin/bash as default entrypoint.
	return []string{
		"/bin/bash",
		"-c",
	}
}

func NewSpecFileFromYaml(data []byte, f string) (*SpecFile, error) {
	ans := &SpecFile{}
	if err := yaml.Unmarshal(data, ans); err != nil {
		return nil, err
	}

	ans.File = f

	return ans, nil
}
