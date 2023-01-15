package cmd

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/chhetripradeep/chtop/pkg/chtop"
)

var (
	cfgFile              string
	clickhouseMetricsUrl string
	clickhouseQueriesUrl string
	clickhouseDatabase   string
	clickhouseUsername   string
	clickhousePassword   string
)

var rootCmd = &cobra.Command{
	Use:   "chtop",
	Short: "ClickHouse monitoring tool",
	Long:  "Monitor your ClickHouse clusters without ever leaving your terminal",
	Run: func(cmd *cobra.Command, _ []string) {
		err := chtop.Run(clickhouseMetricsUrl, clickhouseQueriesUrl, clickhouseDatabase, clickhouseUsername, clickhousePassword)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "path of the config file (default: $HOME/.chtop.yaml)")
	rootCmd.PersistentFlags().StringVarP(&clickhouseMetricsUrl, "metrics-url", "m", "", "clickhouse url for pulling metrics in prometheus exposition format")
	rootCmd.PersistentFlags().StringVarP(&clickhouseQueriesUrl, "queries-url", "q", "", "clickhouse endpoint for running clickhouse queries via native protocol")
	rootCmd.PersistentFlags().StringVarP(&clickhouseDatabase, "queries-database", "d", "system", "clickhouse database for connecting from clickhouse client")
	rootCmd.PersistentFlags().StringVarP(&clickhouseUsername, "queries-username", "u", "default", "clickhouse username for running clickhouse queries")
	rootCmd.PersistentFlags().StringVarP(&clickhousePassword, "queries-password", "p", "", "clickhouse password of the provided clickhouse user for running clickhouse queries")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		cobra.CheckErr(err)
		viper.AddConfigPath(home)
		viper.SetConfigName(".chtop")
	}

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err == nil {
		fmt.Fprintln(os.Stderr, "using config file:", viper.ConfigFileUsed())
	}
}
