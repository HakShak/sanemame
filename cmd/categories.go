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

// categoriesCmd represents the categories command
var categoriesCmd = &cobra.Command{
	Use:   "categories",
	Short: "List categories from Catver.ini",
	Long:  `Load and list the set of categories from Catver.ini`,
	Run: func(cmd *cobra.Command, args []string) {
		dbPath := viper.GetString(DatabaseLocation)
		boltDb, err := bolt.Open(dbPath, 0600, nil)
		if err != nil {
			log.Fatal(err)
		}
		defer boltDb.Close()

		categories := db.GetRawCategories(boltDb)

		tw := new(tabwriter.Writer)
		tw.Init(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(tw, "Category")
		fmt.Fprintln(tw, "--------")

		for _, category := range categories {
			fmt.Fprintf(tw, "%s\n", category)
		}

		fmt.Fprintln(tw)
		tw.Flush()
	},
}

var categoriesStatCmd = &cobra.Command{
	Use:   "categories",
	Short: "Stat Catver.ini",
	Long:  `Load stats from Catver.ini`,
	Run: func(cmd *cobra.Command, args []string) {
		dbPath := viper.GetString(DatabaseLocation)
		boltDb, err := bolt.Open(dbPath, 0600, nil)
		if err != nil {
			log.Fatal(err)
		}
		defer boltDb.Close()

		tw := new(tabwriter.Writer)
		tw.Init(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(tw, "Type\tCount\tMachines")
		fmt.Fprintln(tw, "----\t-----\t--------")

		rawCategories := db.GetRawCategoryMachines(boltDb)
		rawMachines := db.UniqueStrings(rawCategories)
		fmt.Fprintf(tw, "Raw\t%d\t%d\n", len(rawCategories), len(rawMachines))

		primaryCategories := db.GetPrimaryCategoryMachines(boltDb)
		primaryMachines := db.UniqueStrings(primaryCategories)
		fmt.Fprintf(tw, "Primary\t%d\t%d\n", len(primaryCategories), len(primaryMachines))

		secondaryCategories := db.GetSecondaryCategoryMachines(boltDb)
		secondaryMachines := db.UniqueStrings(secondaryCategories)
		fmt.Fprintf(tw, "Secondary\t%d\t%d\n", len(secondaryCategories), len(secondaryMachines))

		matureMachines := rawCategories[db.MatureCategory]
		fmt.Fprintf(tw, "Mature\t%d\t%d\n", 1, len(matureMachines))

		fmt.Fprintln(tw)
		tw.Flush()
	},
}

func init() {
	listCmd.AddCommand(categoriesCmd)
	statCmd.AddCommand(categoriesStatCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// categoriesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// categoriesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
