package main

import (
	"fmt"
	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"okta-aws-role-selector/saml"
	"os"

	"net/http"
	"okta-aws-role-selector/handlers"
	"strings"
)

const (
	ConfigPrefix = "selector"
	ConfigType   = "yaml"

	DefaultTemplatePage = "welcome"
	RolesTemplatePage   = "roles"
	ExampleTemplatePage = "example"
)

var (
	cfgFile         string
	configExtension = fmt.Sprintf(".%s", ConfigType)
)

var serverCmd = &cobra.Command{
	Use:   "Okta-AWS-Role-Selector",
	Short: "Interceptor app for selecting an AWS role from Okta",
	RunE: func(cmd *cobra.Command, args []string) error {
		// load config
		serverConfig := new(saml.Config)
		allSettings := viper.AllSettings()["server"]
		err := mapstructure.Decode(allSettings, serverConfig)
		if err != nil {
			log.Printf("Failed to load server config: %s", err)
			return err
		}

		// setup API handlers
		api := gin.Default()
		api.HTMLRender = ginview.Default()
		api.Use(static.Serve("/", static.LocalFile("assets/", false)))

		api.GET("/", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, DefaultTemplatePage, gin.H{})
		})

		roleHandler, err := handlers.RolesHandler(RolesTemplatePage, serverConfig)
		if err != nil {
			log.Printf("Failed to initialize roles handler: %s", err)
			return err
		}
		api.POST("/", roleHandler)

		exampleHandler, err := handlers.ExampleHandler(ExampleTemplatePage, serverConfig)
		if err != nil {
			log.Printf("Failed to initialize example handler: %s", err)
			return err
		}
		api.POST("/example", exampleHandler)

		return api.Run(":80")
	},
}

func main() {
	if err := serverCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	serverCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "config.yaml", "Config file for Saasy API server")
}

func initConfig() {
	// environment vars integration
	viper.AutomaticEnv()
	viper.SetEnvPrefix(ConfigPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// config setup
	viper.SetConfigFile(cfgFile)
	viper.SetConfigType(ConfigType)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Printf("Not found: %s\n", err)
			os.Exit(1)
		} else {
			log.Printf("Failed to read config file: %s\n", err)
			os.Exit(1)
		}
	}

	// set log level based on config
	logLevelRaw := viper.GetString("log_level")
	logLevel, err := log.ParseLevel(logLevelRaw)
	if err != nil {
		log.Printf("Failed to parse log level: %s\n", err)
		os.Exit(1)
	}
	log.SetLevel(logLevel)
	log.Printf("Using config file: %s\n", viper.ConfigFileUsed())
}
