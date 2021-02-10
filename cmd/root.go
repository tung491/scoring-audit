package cmd

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var (
	cfgFile  string
	userName string
	token    string
)

var rootCmd = &cobra.Command{
	Use:   "audit",
	Short: "Audit tasks which missing components",
	Long:  "Audit tasks which missing components",
	PreRun: func(cmd *cobra.Command, args []string) {
		fmt.Println("Pre run")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.jira.yaml)")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigName(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		viper.AddConfigPath(home)
		viper.SetConfigName(".jira")
	}
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Using config file ", viper.ConfigFileUsed())
	}
	userName = viper.GetString("username")
	token = viper.GetString("token")
}
