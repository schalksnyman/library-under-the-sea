package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"library-under-the-sea/services/library-repo/internal/agent"
)

func main() {
	cli := &cli{}

	cmd := &cobra.Command{
		Use:     "library-repo",
		PreRunE: cli.setupConfig,
		RunE:    cli.run,
	}

	if err := setupFlags(cmd); err != nil {
		log.Fatal(err)
	}

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

type cli struct {
	cfg cfg
}

type cfg struct {
	agent.Config
}

func setupFlags(cmd *cobra.Command) error {
	cmd.Flags().String("config-file", "", "Path to config file.")
	cmd.Flags().String("http-addr",
		"127.0.0.1:8401",
		"Service library-repo healthcheck address")
	cmd.Flags().String("grpc-addr",
		"127.0.0.1:8400",
		"Service grpc server address")
	cmd.Flags().String("db-connect-string",
		"mongodb://localhost:27017",
		"Database connect string")
	cmd.Flags().String("db-name",
		"library",
		"Database name")

	return viper.BindPFlags(cmd.Flags())
}

func (c *cli) setupConfig(cmd *cobra.Command, args []string) error {
	var err error

	configFile, err := cmd.Flags().GetString("config-file")
	if err != nil {
		return err
	}
	viper.SetConfigFile(configFile)

	if err = viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	c.cfg.HTTPAddr = viper.GetString("http-addr")
	c.cfg.GRPCAddr = viper.GetString("grpc-addr")
	c.cfg.DBConnectString = viper.GetString("db-connect-string")
	c.cfg.DBName = viper.GetString("db-name")

	return nil
}

func (c *cli) run(cmd *cobra.Command, args []string) error {
	var err error
	a, err := agent.New(c.cfg.Config)
	if err != nil {
		return err
	}
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	<-sigc
	return a.Shutdown()
}
