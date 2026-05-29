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
	"regexp"

	loader "github.com/macaroni-os/whip/pkg/loader"
	specs "github.com/macaroni-os/whip/pkg/specs"

	tablewriter "github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
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
			filter, _ := cmd.Flags().GetString("filter")
			var r *regexp.Regexp

			config.AddSpecsDirs(specsDir)

			// Check instance
			whip := loader.NewWhipHolder(config)

			// Load rules
			err := whip.LoadSpecs()
			if err != nil {
				fmt.Println("Error on load rules: " + err.Error())
				os.Exit(1)
			}

			if filter != "" {
				r = regexp.MustCompile(filter)
				if r == nil {
					fmt.Println("Invalid regex: %s", filter)
					os.Exit(1)
				}
			}

			if tableMode {
				table := tablewriter.NewTable(os.Stdout,
					tablewriter.WithRenderer(renderer.NewBlueprint(tw.Rendition{
						Settings: tw.Settings{
							Separators: tw.Separators{
								BetweenRows: tw.On,
							},
						},
					})))
				table.Header("File", "Hook", "Description")

				for f, v := range *whip.GetRules() {

					for n, h := range v.Hooks {
						hName := fmt.Sprintf("%s.%s", f, n)
						if r == nil || (r != nil && r.MatchString(hName)) {
							table.Append([]string{
								f, hName, h.Description,
							})
						}
					}
				}

				table.Render()

			} else {

				for f, v := range *whip.GetRules() {
					for n, _ := range v.Hooks {
						hName := fmt.Sprintf("%s.%s", f, n)
						if r == nil || (r != nil && r.MatchString(hName)) {
							fmt.Println(fmt.Sprintf("%s - %s.%s",
								f, f, n))
						}
					}
				}
			}

		},
	}

	flags := cmd.Flags()
	flags.StringSlice("specs-dir", []string{},
		"Define additional specs directory at runtime.")
	flags.Bool("table", false, "Show hooks in table format.")
	flags.StringP("filter", "f", "", "Filter hook matching string.")

	return cmd
}
