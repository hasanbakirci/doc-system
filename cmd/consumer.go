/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"github.com/hasanbakirci/doc-system/cmd/listener"
	"github.com/hasanbakirci/doc-system/internal/config"

	"github.com/spf13/cobra"
)

// consumerCmd represents the consumer command
var consumerCmd = &cobra.Command{
	Use:   "consumer",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("consumer called")
	},
}

func init() {
	rootCmd.AddCommand(consumerCmd)

	var cfgFile string
	consumerCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "config.dev", "config file (default is $HOME/.golang-api.yaml)")

	ApiConfig, err := config.GetAllValues("./config/", cfgFile)
	if err != nil {
		panic(err)
	}

	consumerCmd.Run = func(cmd *cobra.Command, args []string) {

		forever := make(chan bool)

		l := listener.NewListener(*ApiConfig)
		l.Start()

		<-forever

	}
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// consumerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// consumerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
