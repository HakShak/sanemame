package cmd

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

import "github.com/HakShak/sanemame/mamexml"

// controlsCmd represents the controls command
var controlsCmd = &cobra.Command{
	Use:   "controls",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		controls, err := mamexml.LoadControlsXml("controls.xml")
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Controls: %d", len(controls))

	},
}

var listControlsCmd = &cobra.Command{
	Use:   "controls",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		controls, err := mamexml.LoadControlsXml("controls.xml")
		if err != nil {
			log.Fatal(err)
		}

		controlList := make(map[string]map[string]bool)

		for _, game := range controls {
			for _, player := range game.PlayerControls {
				for _, control := range player.Controls {
					if _, ok := controlList[control.Constant.Name]; !ok {
						controlList[control.Constant.Name] = make(map[string]bool)
					}
					controlList[control.Constant.Name][control.Name] = true
				}
			}
		}

		tw := new(tabwriter.Writer)
		tw.Init(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(tw, "Keyword\tDescription")
		fmt.Fprintln(tw, "-------\t-----------")

		for keyword, descriptions := range controlList {
			for description, _ := range descriptions {
				fmt.Fprintf(tw, "%s\t%v\n", keyword, description)
			}
		}

		fmt.Fprintln(tw)
		tw.Flush()
	},
}

func init() {
	statCmd.AddCommand(controlsCmd)
	listCmd.AddCommand(listControlsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// controlsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// controlsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
