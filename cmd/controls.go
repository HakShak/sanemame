package cmd

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/HakShak/sanemame/db"
	"github.com/boltdb/bolt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

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
		dbPath := viper.GetString(DatabaseLocation)
		boltDb, err := bolt.Open(dbPath, 0600, nil)
		if err != nil {
			log.Fatal(err)
		}
		defer boltDb.Close()

		controls := db.GetControlMachines(boltDb)
		controlMachines := db.UniqueStrings(controls)

		tw := new(tabwriter.Writer)
		tw.Init(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(tw, "Type\tCount\tMachines")
		fmt.Fprintln(tw, "----\t-----\t--------")

		fmt.Fprintf(tw, "%s\t%d\t%d\n", "Controls", len(controls), len(controlMachines))

		fmt.Fprintln(tw)
		tw.Flush()

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
		dbPath := viper.GetString(DatabaseLocation)
		boltDb, err := bolt.Open(dbPath, 0600, nil)
		if err != nil {
			log.Fatal(err)
		}
		defer boltDb.Close()

		controls := db.GetControlNames(boltDb)

		tw := new(tabwriter.Writer)
		tw.Init(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(tw, "Keyword\tNames")
		fmt.Fprintln(tw, "-------\t-----")

		for control, desc := range controls {
			fmt.Fprintf(tw, "%s\t%s\n", control, desc)
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
