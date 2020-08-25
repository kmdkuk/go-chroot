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
	"os"
	"os/exec"
	"syscall"

	"github.com/kmdkuk/go-chroot/log"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"

	"path/filepath"

	"github.com/docker/docker/pkg/reexec"
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

	reexec.Register("nsInit", nsInit)
	if reexec.Init() {
		log.Fatal("reexec.Init() error")
	}
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
	targetDir := filepath.Join(os.TempDir(), "go-chroot")
	defer func() {
		err = os.RemoveAll(targetDir)
		if err != nil {
			log.Fatal(err)
		}
	}()
	err = chrootExecSh(targetDir)
	if err != nil {
		log.Fatal(err)
	}
}

func extractToTemp(filename string) error {
	targetDir := filepath.Join(os.TempDir(), "go-chroot")
	f, err := os.Stat(targetDir)
	if os.IsNotExist(err) || f.IsDir() {
		err = os.MkdirAll(targetDir, 0755)
		if err != nil {
			return err
		}
	}
	err = exec.Command("tar", "-xvf", filename, "-C", targetDir).Run()
	if err != nil {
		return err
	}
	return nil
}

func chrootExecSh(targetDir string) error {
	// err := syscall.Chroot(targetDir)
	// if err != nil {
	// 	return err
	// }
	// err = syscall.Chdir("/")
	// if err != nil {
	// 	return err
	// }
	cmd := reexec.Command("nsInit")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Env = []string{"PS1=-[ns-process]- # "}

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWNS |
			syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWIPC |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNET |
			syscall.CLONE_NEWUSER,
		UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Geteuid(),
				Size:        1,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getgid(),
				Size:        1,
			},
		},
	}

	return cmd.Run()
}

func nsInit() {
	log.Debug("namespace setup code goes here")
	nsRun()
}

func nsRun() {
	cmd := exec.Command("/bin/sh")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Env = []string{"PS1=-[ns-process]- # "}

	if err := cmd.Run(); err != nil {
		log.Fatalf("Error runnning the /bin/sh command - %s\n", err)
	}
}
