package main

import (
	"errors"
	"os"
	"os/signal"
	"porage/internal/pkg"
	"porage/internal/server"

	"github.com/spf13/cobra"
)

var (
	configFilePath string
)

func main() {
	err := parseCommandLine()
	if err != nil {
		panic(err)
	}

	config, err := pkg.ParseConfigFile(configFilePath)
	if err != nil {
		panic(err)
	}

	server := server.NewPorageServer(config)

	// Wait for SIGINT (Ctrl+C) to stop the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		server.Stop()
	}()

	server.Start()
}

func parseCommandLine() error {
	var rootCmd = &cobra.Command{
		Use:   "porage-server",
		Short: "Porage server application",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Check if the configFilePath is empty
			if configFilePath == "" {
				return errors.New("configuration file path is required but not provided")
			}
			return nil
		},
	}

	// Adding a flag for the configuration file
	rootCmd.PersistentFlags().StringVarP(&configFilePath, "config", "c", "", "Config file (default is ./config.toml)")

	// Execute the command
	if err := rootCmd.Execute(); err != nil {
		return err
	}

	return nil
}
