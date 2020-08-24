/*
Copyright © 2020 kouki kamada(kmdkuk.com)

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"archive/tar"
	"fmt"
	"github.com/kmdkuk/go-chroot/log"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-chroot",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: Run,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.go-chroot.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			log.Fatal(err)
		}

		// Search config in home directory with name ".go-chroot" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".go-chroot")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Error("Using config file:", viper.ConfigFileUsed())
	}
}

func Run(_ *cobra.Command, _ []string) {
	log.Debug(os.TempDir())
	err := extractToTemp("alpine.tar")
	if err != nil {
		log.Fatal(err)
	}
	err = os.RemoveAll(filepath.Join(os.TempDir(),"go-chroot"))
	if err != nil {
		log.Fatal(err)
	}
}

func extractToTemp(filename string) error {
	e := filepath.Ext(filename)
	if e != ".tar" {
		return fmt.Errorf("tar以外の拡張子に対応してません．: %s", e)
	}
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	targetDir := filepath.Join(os.TempDir(), "go-chroot")
	reader := tar.NewReader(file)
	for {
		header, err := reader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		target := path.Join(targetDir, header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			err := os.MkdirAll(target, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			setAttrs(target, header)
			break
		case tar.TypeReg, tar.TypeSymlink:
			w, err := os.Create(target)
			if err != nil{
				return err
			}
			_, err = io.Copy(w, reader)
			if err != nil {
				return err
			}
			w.Close()

			setAttrs(target, header)
			break
		default:
			return fmt.Errorf("unsupported type: %v", string(header.Typeflag))
			break
		}
	}
	return nil
}

func setAttrs(target string, header *tar.Header){
	os.Chmod(target, os.FileMode(header.Mode))
	os.Chtimes(target, header.AccessTime, header.ModTime)
}
