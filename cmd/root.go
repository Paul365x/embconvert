/*
Copyright © 2026 Paul Chubb <paulc@singlespoon.org>
*/
package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "embconvert",
	Short: "does actions on a tree of files",
	Long: `Walks a directory tree and does actions on found matching files:
	- copy matching file to outpath
	- create metadata pdf from file
	- create json metadata files
	- count matching files`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.AddCommand(reportCmd)
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.embconvert.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

/*
** Commands
*/

// counts instances of this type of file in the give tree.
var reportCmd = &cobra.Command{
	Use:   "report <inpath> <filetype>",
	Short: "Count the number of this filetype in the inpath tree",
	Args:  cobra.MaximumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 { 
			return fmt.Errorf("Need to specify both <inpath> and <filetype>")    
		}
		path := args[0]
		ext := args[1]
		if !strings.HasPrefix(ext,".") {
			ext = "." + ext
		}
		_, err := os.Stat(path)
		if err != nil {
			return err 
		}
		
		// setup the action
		var fileCount = 0;
		act := func (string) {
    		fileCount++
		}

		// setup the id function
		id := func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			tgtExt := filepath.Ext(d.Name())
			if ext == tgtExt {
				act(path)
			}
			//fmt.Printf("%s %s : %s\n", d.Name(), ext, tgtExt)
			return nil

		}
		// Walk the tree
		err = filepath.WalkDir(path,id)
		if err != nil {
			return err
		}
		fmt.Printf("Found %d %s files in %s\n", fileCount, ext, path)
		return nil
	},
}


/*
** Utility functions
*/

type IdFn func(string,string) error
type ActionFn func(string) bool

