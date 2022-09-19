/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"time"

	"github.com/hasanbakirci/doc-system/internal/auth"
	"github.com/hasanbakirci/doc-system/internal/config"
	"github.com/hasanbakirci/doc-system/internal/document"
	elasticclient "github.com/hasanbakirci/doc-system/pkg/elasticClient"
	"github.com/hasanbakirci/doc-system/pkg/graceful"
	"github.com/hasanbakirci/doc-system/pkg/middleware"
	"github.com/hasanbakirci/doc-system/pkg/redisClient"
	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
)

// apiCmd represents the api command
var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	rootCmd.AddCommand(apiCmd)

	var port string
	var cfgFile string
	apiCmd.PersistentFlags().StringVarP(&port, "port", "p", "9494", "Restfull Service Port")
	apiCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "config.dev", "config file (default is $HOME/.golang-api.yaml)")

	ApiConfig, err := config.GetAllValues("./config/", cfgFile)
	if err != nil {
		panic(err)
	}
	apiCmd.Run = func(cmd *cobra.Command, args []string) {
		instance := echo.New()

		instance.Use(middleware.RecoveryMiddlewareFunc, middleware.LoggingMiddlewareFunc)

		// db, err := mongoClient.ConnectDb(ApiConfig.MongoSettings)
		// if err != nil {
		// 	fmt.Println("Db connection error")
		// }

		redis := redisClient.NewRedisClient(ApiConfig.RedisSettings.Uri)

		elastic, err := elasticclient.ConnectElastic()

		if err != nil {
			fmt.Println("Elastic connection error")
		}
		// document
		//documentRepository := document.NewDocumentRepository(db)
		documentRepository := document.NewElasticRepository(elastic)
		documentService := document.NewDocumentService(documentRepository, redis)
		documentHandler := document.NewDocumentHandler(documentService)
		document.RegisterDocumentHandlers(instance, documentHandler, ApiConfig.JwtSettings.SecretKey)
		// auth
		//authRepository := auth.NewAuthRepository(db)
		authRepository := auth.NewElasticRepository(elastic)
		authService := auth.NewAuthService(authRepository, *ApiConfig)
		authHandler := auth.NewUserHandler(authService)
		auth.RegisterUserHandlers(instance, authHandler)

		fmt.Println("Api starting")
		if err := instance.Start(fmt.Sprintf(":%s", port)); err != nil {
			fmt.Println("Api fatal error")
		}
		graceful.Shutdown(instance, time.Second*2)
	}
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// apiCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// apiCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
