package cmd

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

import "github.com/HakShak/sanemame/mamexml"

// nplayersCmd represents the nplayers command
var nplayersCmd = &cobra.Command{
	Use:   "nplayers",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		nplayers, err := mamexml.LoadNPlayersIni("nplayers.ini")
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("NPlayers: %d", len(nplayers))

	},
}

var listNplayersCmd = &cobra.Command{
	Use:   "nplayers",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		nplayers, err := mamexml.LoadNPlayersIni("nplayers.ini")
		if err != nil {
			log.Fatal(err)
		}

		playerTypes := make(map[string]bool)

		for _, players := range nplayers {
			for _, player := range players {
				playerTypes[player.Raw] = true
			}
		}

		tw := new(tabwriter.Writer)
		tw.Init(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(tw, "Keyword")
		fmt.Fprintln(tw, "-------")

		for players := range playerTypes {
			fmt.Fprintf(tw, "%s\n", players)
		}

		fmt.Fprintln(tw)
		tw.Flush()

	},
}

func init() {
	statCmd.AddCommand(nplayersCmd)
	listCmd.AddCommand(listNplayersCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// nplayersCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// nplayersCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
