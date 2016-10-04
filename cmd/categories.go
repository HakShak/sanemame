package cmd

import (
	"fmt"
	"github.com/HakShak/sanemame/mamexml"
	"github.com/spf13/cobra"
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
		categories, err := mamexml.LoadCatverIni("Catver.ini")
		if err != nil {
			log.Fatal(err)
		}

		tw := new(tabwriter.Writer)
		tw.Init(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(tw, "Primary\tSecondary")
		fmt.Fprintln(tw, "-------\t---------")

		categorySet := make(map[string]map[string]bool)

		for _, category := range categories {
			if _, ok := categorySet[category.Primary]; !ok {
				categorySet[category.Primary] = make(map[string]bool)
			}
			categorySet[category.Primary][category.Secondary] = category.Mature
		}

		for primary, secondarySet := range categorySet {
			for secondary, _ := range secondarySet {
				if secondary != "" {
					fmt.Fprintf(tw, "%s\t%s\n", primary, secondary)
				}
			}
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
