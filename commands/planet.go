package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var RootCmd = &cobra.Command{
	Use:   "Bastion",
	Short: "Bastion is an awesome planet style RSS aggregator",
	Long: `Bastion provides planet style RSS aggregation.
	It is inspired by python planet and has a simple YAML configuration
	and provides it's own web server`,
	Run: rootRun,
}

func rootRun(cmd *cobra.Command, args []string) {
	fmt.Println(viper.Get("feeds"))
	fmt.Println(viper.GetString("appname"))
}


func addCommands() {
	RootCmd.AddCommand(fetchCmd)
}

func Execute() {
	addCommands()
	err := RootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

var CfgFile string

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&CfgFile, "config", "", "config file (default is "+
		"$HOME/.bastion/config.yaml")
}

func initConfig() {
	if CfgFile != "" {
		viper.SetConfigFile(CfgFile)
	}
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/bastion/")
	viper.AddConfigPath("$HOME/.bastion/")
	viper.ReadInConfig()
}
