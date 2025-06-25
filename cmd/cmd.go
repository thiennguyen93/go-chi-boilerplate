package cmd

import (
	"fmt"
	"log"
	"os"

	"thiennguyen.dev/welab-healthcare-app/config"
	"thiennguyen.dev/welab-healthcare-app/routers"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

const defaultConfigFilename = "config.deployment.yml"
const defaultRuntimeEnv = "development"

var serverConfig config.Configuration

// New returns a new CLI app.
// Build with:
// $ go build -ldflags="-s -w"
func New() *cobra.Command {
	configFile := defaultConfigFilename
	runtimeEnv := defaultRuntimeEnv

	rootCmd := &cobra.Command{
		Use:                        "project",
		Short:                      "Command line interface for project.",
		Long:                       "The root command will start the HTTP server.",
		Version:                    "v0.0.1",
		SilenceErrors:              true,
		SilenceUsage:               true,
		TraverseChildren:           true,
		SuggestionsMinimumDistance: 1,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Read configuration from file before any of the commands run functions.
			if runtimeEnv == "development" || runtimeEnv == "local" || runtimeEnv == "deployment" {
				envFile := fmt.Sprintf(".env.%s", runtimeEnv)
				err := godotenv.Load(envFile) // 👈 load .env.{runtimeEnv} file
				fmt.Println("LOADED ", envFile)
				if err != nil {
					// Missing .env
					log.Fatal(err)
				}
				os.Setenv("RUNTIME_ENV", runtimeEnv)
			} else {
				os.Setenv("RUNTIME_ENV", "production")
			}
			fmt.Println("Environment =", os.Getenv("RUNTIME_ENV"))
			return serverConfig.BindFile(configFile)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return startServer()
		},
	}

	// Shared flags.
	flags := rootCmd.PersistentFlags()
	// Config file flag
	flags.StringVar(&configFile, "config", configFile, "--config=config.deployment.yml a filepath which contains the YAML config format")
	// Runtime environment flag
	flags.StringVar(&runtimeEnv, "env", runtimeEnv, "--env=development specify runtime environment")

	// Subcommands here.
	// Subcommand: Make migration
	subCommands := []cobra.Command{
		{
			Use:   "seed",
			Short: "Seed the database with sample data",
			RunE: func(cmd *cobra.Command, args []string) error {
				fmt.Println("Start seeding")
				// seeding.Seed(serverConfig)
				return nil
			},
		},
		{
			Use:   "migrate",
			Short: "Migrate the database structure",
			RunE: func(cmd *cobra.Command, args []string) error {
				fmt.Println("Start migrate")
				// migrate.Migrate(serverConfig)
				return nil
			},
		},
		{
			Use:   "goose",
			Short: "Migration DB with goose",
			RunE: func(cmd *cobra.Command, args []string) error {
				// goose.Migrate(serverConfig, args)
				return nil
			},
		},
	}

	for i := 0; i < len(subCommands); i++ {
		rootCmd.AddCommand(&subCommands[i])
	}

	return rootCmd
}

func startServer() error {
	srv := config.NewServer(serverConfig)
	routers.RegisterRoutes(srv.Echo)
	return srv.Start()
}
