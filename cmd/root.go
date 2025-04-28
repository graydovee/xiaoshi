package cmd

import (
	"github.com/graydovee/xiaoshi/pkg/config"
	"github.com/spf13/cobra"
	"log/slog"
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
		slog.SetLogLoggerLevel(slog.Level(logLevel))
		cfg, err := config.LoadConfigFromPath()
		if err != nil {
			slog.Error("load config error: ", err)
			return
		}

		Run(cfg)
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

var logLevel int

func init() {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	rootCmd.Flags().IntVarP(&logLevel, "log-level", "l", int(slog.LevelDebug), "config file path")
}
