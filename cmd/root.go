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

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default: $HOME/.chtop.yaml)")
	rootCmd.PersistentFlags().StringVar(&clickhouseMetricsUrl, "metrics-url", "", "clickhouse url for pulling metrics in prometheus exposition format")
	rootCmd.PersistentFlags().StringVar(&clickhouseQueriesUrl, "queries-url", "", "clickhouse url for running clickhouse queries (native protocol port)")
	rootCmd.PersistentFlags().StringVar(&clickhouseDatabase, "queries-database", "system", "clickhouse database for connecting clickhouse client")
	rootCmd.PersistentFlags().StringVar(&clickhouseUsername, "queries-username", "default", "clickhouse username for running clickhouse queries")
	rootCmd.PersistentFlags().StringVar(&clickhousePassword, "queries-password", "", "clickhouse password for running clickhouse queries")
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
