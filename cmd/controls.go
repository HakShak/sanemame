package cmd

import (
	"log"

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

func init() {
	statCmd.AddCommand(controlsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// controlsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// controlsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
