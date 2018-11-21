package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "go-assets",
	Short: "create an assets package reflecting the contents of a single folder",
	Long: `
	
	you can do anything at zombocom

	`,
}

func Execute() {
	log.SetFlags(0)
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
