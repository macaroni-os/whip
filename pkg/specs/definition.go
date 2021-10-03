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

type SpecFile struct {
	File string `yaml:"-" json:"-"`

	Hooks      map[string]Hook `yaml:"hooks" json:"hooks"`
	Entrypoint []string        `yaml:"entrypoint,omitempty" json:"entrypoint,omitempty"`
}

type Hook struct {
	Remediate   string   `yaml:"remediate,omitempty" json:"remediate,omitempty"`
	Check       string   `yaml:"check,omitempty" json:"check,omitempty"`
	Description string   `yaml:"description,omitempty" json:"description,omitempty"`
	Keywords    []string `yaml:"keywords,omitempty" json:"keywords,omitempty"`
	Actions     []string `yaml:"actions.omitempty" json:"actions,omitempty"`
}
