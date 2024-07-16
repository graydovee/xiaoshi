/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"git.graydove.cn/graydove/xiaoshi.git/pkg/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	easy "github.com/t-tomalak/logrus-easy-formatter"
	"gopkg.in/yaml.v3"
	"os"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "run",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		log.Info("loading config: ", configFile)

		cfgbin, err := os.ReadFile(configFile)
		if err != nil {
			panic(err)
		}
		var cfg config.Config
		if err = yaml.Unmarshal(cfgbin, &cfg); err != nil {
			panic(err)
		}

		Run(&cfg)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var configFile string
var logLevel int

func init() {
	log.SetFormatter(&easy.Formatter{
		TimestampFormat: "2006-01-02 15:04:05",
		LogFormat:       "[zero][%time%][%lvl%]: %msg% \n",
	})
	log.SetLevel(log.DebugLevel)

	rootCmd.Flags().StringVarP(&configFile, "config", "c", "config.yaml", "config file path")
	rootCmd.Flags().IntVarP(&logLevel, "log-level", "l", int(log.DebugLevel), "config file path")
}
