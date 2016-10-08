package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

import "github.com/HakShak/sanemame/mamexml"

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
		api := viper.GetString(GithubReleasesApi)
		repo := viper.GetString(MameRepo)
		filename, err := mamexml.GetLatestXmlFile(api, repo)
		if err != nil {
			log.Fatal(err)
		}

		machines, err := mamexml.Load(filename)
		if err != nil {
			log.Fatal(err)
		}

		devices := 0
		bios := 0
		runnable := 0
		mechanical := 0
		clones := 0
		roms := 0
		samples := 0

		for _, m := range machines {
			if m.IsDevice {
				devices += 1
			}
			if m.IsBios {
				bios += 1
			}
			if m.IsRunnable {
				runnable += 1
			}
			if m.IsMechanical {
				mechanical += 1
			}
			if m.CloneOf != "" {
				clones += 1
			}
			if m.RomOf != "" {
				roms += 1
			}
			if m.SampleOf != "" {
				samples += 1
			}
		}

		log.Printf("Machines: %d", len(machines))
		log.Printf("Devices: %d", devices)
		log.Printf("Bios: %d", bios)
		log.Printf("Runnable: %d", runnable)
		log.Printf("Mechanical: %d", mechanical)
		log.Printf("Clones: %d", clones)
		log.Printf("Samples: %d", samples)
		log.Printf("Roms: %d", roms)
		potential := len(machines) - devices - bios - mechanical
		log.Printf("Potential: %d", potential)
		nonClones := runnable - clones
		log.Printf("NonClones: %d", nonClones)
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
