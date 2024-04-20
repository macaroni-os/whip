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

	loader "github.com/macaroni-os/whip/pkg/loader"
	specs "github.com/macaroni-os/whip/pkg/specs"

	tablewriter "github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func newListCommand(config *specs.Config) *cobra.Command {
	var cmd = &cobra.Command{
		Use:     "list",
		Short:   "List available hooks.",
		Aliases: []string{"h"},
		PreRun: func(cmd *cobra.Command, args []string) {
		},
		Run: func(cmd *cobra.Command, args []string) {

			specsDir, _ := cmd.Flags().GetStringSlice("specs-dir")
			tableMode, _ := cmd.Flags().GetBool("table")

			config.AddSpecsDirs(specsDir)

			// Check instance
			whip := loader.NewWhipHolder(config)

			// Load rules
			err := whip.LoadSpecs()
			if err != nil {
				fmt.Println("Error on load rules: " + err.Error())
				os.Exit(1)
			}

			if tableMode {
				table := tablewriter.NewWriter(os.Stdout)
				table.SetBorders(tablewriter.Border{
					Left:   true,
					Top:    true,
					Right:  true,
					Bottom: true,
				})
				table.SetHeader([]string{
					"File", "Name", "Description"})
				table.SetAutoWrapText(false)

				for f, v := range *whip.GetRules() {
					for n, h := range v.Hooks {
						table.Append([]string{f, n, h.Description})
					}
				}

				table.Render()

			} else {

				for f, v := range *whip.GetRules() {
					for n, _ := range v.Hooks {
						fmt.Println(fmt.Sprintf("%s - %s",
							f, n))
					}
				}
			}

		},
	}

	flags := cmd.Flags()
	flags.StringSlice("specs-dir", []string{},
		"Define additional specs directory at runtime.")
	flags.Bool("table", false, "Show hooks in table format.")

	return cmd
}
