package cmd

import (
	"fmt"
	"github.com/HakShak/sanemame/db"
	"github.com/HakShak/sanemame/mamexml"
	"github.com/boltdb/bolt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"text/tabwriter"
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

		categories := db.GetCategories(boltDb)

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
		categories, err := mamexml.LoadCatverIni("Catver.ini")
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Categorized: %d", len(categories))
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
