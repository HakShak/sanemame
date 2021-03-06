package cmd

import (
	"log"
	"os"

	"github.com/dustin/go-humanize"

	"github.com/HakShak/sanemame/db"
	"github.com/HakShak/sanemame/filetypes/mamexml"
	"github.com/boltdb/bolt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// initdbCmd represents the initdb command
var initdbCmd = &cobra.Command{
	Use:   "initdb",
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

		db.UpdateCategories(boltDb, "Catver.ini")
		db.UpdateControls(boltDb, "controls.xml")
		db.UpdateNPlayers(boltDb, "nplayers.ini")

		fileName, err := mamexml.GetLatestXMLFile(
			viper.GetString(GithubReleasesAPI),
			viper.GetString(MameRepo))
		if err != nil {
			log.Fatal(err)
		}

		db.UpdateMachines(boltDb, fileName)

		fileInfo, err := os.Stat(dbPath)
		if err != nil {
			log.Fatal(err)
		}

		humanSize := humanize.Bytes(uint64(fileInfo.Size()))

		log.Printf("Database Size: %s", humanSize)
	},
}

func init() {
	RootCmd.AddCommand(initdbCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initdbCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initdbCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
