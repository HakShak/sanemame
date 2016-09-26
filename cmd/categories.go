package cmd

import (
	"github.com/HakShak/sanemame/mamexml"
	"github.com/spf13/cobra"
	"log"
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

		log.Printf("Categorized: %d", len(categories))

		categorySet := make(map[string]map[string]bool)

		for _, category := range categories {
			if _, ok := categorySet[category.Primary]; !ok {
				categorySet[category.Primary] = make(map[string]bool)
			}
			categorySet[category.Primary][category.Secondary] = category.Mature
		}

		for primary, secondarySet := range categorySet {
			log.Println(primary)
			for secondary, _ := range secondarySet {
				if secondary != "" {
					log.Printf("\t%s", secondary)
				}
			}
		}
	},
}

func init() {
	listCmd.AddCommand(categoriesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// categoriesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// categoriesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
