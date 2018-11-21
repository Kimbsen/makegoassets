package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"os"
)

var (
	createFlags struct {
		folder string
		prefix string
	}

	createCmd = &cobra.Command{
		Use:     "create",
		Short:   "make the thing",
		Example: "go-assets create -f folder_i_want_to_archive",
		Run: func(cmd *cobra.Command, args []string) {
			cfg := &Config{
				Folder:        createFlags.folder,
				PackagePrefix: createFlags.prefix,
			}
			if err := cfg.Validate(); err != nil {
				log.Println("Invalid state, cannot continue:", err)
				os.Exit(3)
			}
			err := CreatePackage(cfg)
			if err != nil {
				log.Println("Ops. Something went wrong")
				log.Println("Error:", err)
				os.Exit(4)
			}
			log.Printf(`
created and installed package "%s/assets"
		

now just add 

     "%s/assets" 

to your imports and use the assets package
as follows:

   bytes, err := assets.Get("%s/file")

where file is the name of a file in '%s'
subfolders and their files are also added

the only real error is a not found error


to update the package to an updated view of the files
and folders in '%s' run the generated pack.sh

	`, cfg.PackagePrefix, cfg.PackagePrefix, cfg.Folder, cfg.Folder, cfg.Folder)
		},
	}
)

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringVarP(&createFlags.folder, "folder", "f", "", "the name of the folder")
	createCmd.Flags().StringVar(&createFlags.prefix, "packageprefix", "", "the package prefix for the assets package")
}
