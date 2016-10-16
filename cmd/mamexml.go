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

// mamexmlCmd represents the mamexml command
var mamexmlCmd = &cobra.Command{
	Use:   "mamexml",
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

		driverStatus := db.GetDriverStatusMachines(boltDb)
		booleanState := db.GetBooleanMachines(boltDb)

		tw := new(tabwriter.Writer)
		tw.Init(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(tw, "Type\tCount\tMachines")
		fmt.Fprintln(tw, "----\t-----\t--------")

		driverMachines := db.UniqueStrings(driverStatus)
		fmt.Fprintf(tw, "%s\t%d\t%d\n", "Driver Statuses", len(driverStatus), len(driverMachines))

		booleanMachines := db.UniqueStrings(booleanState)
		fmt.Fprintf(tw, "%s\t%d\t%d\n", "Boolean States", len(booleanState), len(booleanMachines))

		fmt.Fprintf(tw, "%s\t%d\t%d\n", "Devices", 1, len(booleanState["IsDevice"]))
		fmt.Fprintf(tw, "%s\t%d\t%d\n", "Bios", 1, len(booleanState["IsBios"]))
		fmt.Fprintf(tw, "%s\t%d\t%d\n", "Runnable", 1, len(booleanState["IsRunnable"]))
		fmt.Fprintf(tw, "%s\t%d\t%d\n", "Mechanical", 1, len(booleanState["IsMechanical"]))
		fmt.Fprintf(tw, "%s\t%d\t%d\n", "Clones", 1, len(booleanState["IsClone"]))
		fmt.Fprintf(tw, "%s\t%d\t%d\n", "Samples", 1, len(booleanState["IsSample"]))

		fmt.Fprintln(tw)
		tw.Flush()
	},
}

func init() {
	statCmd.AddCommand(mamexmlCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mamexmlCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// mamexmlCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
