/*

Copyright (C) 2021  Daniele Rondina <geaaru@sabayonlinux.org>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.:s

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.

*/

package cmd

import (
	"fmt"
	"os"

	loader "github.com/geaaru/whip/pkg/loader"
	specs "github.com/geaaru/whip/pkg/specs"

	"github.com/spf13/cobra"
)

func newHookCommand(config *specs.Config) *cobra.Command {
	var cmd = &cobra.Commad{
		Use:     "hook [filename.hook1] ... [filename.hookN]",
		Short:   "Call a specific hook.",
		Aliases: []string{"h"},
		PreRun: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("No hook selected.")
				os.Exit(1)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {

			// Check instance
			whip := loader.NewWhipHolder(config)

			// Load rules
			err := whip.LoadSpecs()
			if err != nil {
				fmt.Println("Error on load rules: " + err.Error())
				os.Exit(1)
			}

			for _, hook := range args {

				err = whip.RunHook(hook)
				if err != nil {
					fmt.Println(
						fmt.Sprintf("Error on run hook %s: %s",
							hook, err.Error(),
						),
					)
					os.Exit(1)
				}

			}

		},
	}

	return cmd
}
