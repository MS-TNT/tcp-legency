package main

import (
	"demo/tcp/client"
	"demo/tcp/server"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var root = &cobra.Command{
	Use:   "tcp subcommand [flags]",
	Short: "tcp",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		initConfig(cmd)
	},
}

func init() {
	srvCmd := NewServerCommand()
	srvCmd.PersistentFlags().StringP("port", "p", "6789", "tcp server port")
	srvCmd.PersistentFlags().StringP("ip", "i", "localhost", "tcp server ip")
	cliCmd := NewClientCommand()
	cliCmd.PersistentFlags().StringP("port", "p", "6789", "tcp server port")
	cliCmd.PersistentFlags().StringP("host", "a", "127.0.0.1", "tcp server host")

	root.AddCommand(srvCmd)
	root.AddCommand(cliCmd)
	// root.GenBashCompletion(os.Stdout) // enable auto-completion
}

func initConfig(cmd *cobra.Command) {
	v := viper.New()
	v.AutomaticEnv()
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		configName := strings.ReplaceAll(f.Name, "-", "_")
		if !f.Changed && v.IsSet(configName) {
			value := v.Get(configName)
			cmd.Flags().Set(f.Name, fmt.Sprintf("%v", value))
		}
	})
}

func main() {
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

func NewServerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Run tcp server",
		RunE: func(cmd *cobra.Command, args []string) error {
			srv := server.NewServer()
			return srv.StartTcpServer(cmd, args)
		},
	}
	return cmd
}

func NewClientCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "client",
		Short: "Run tcp client",
		RunE: func(cmd *cobra.Command, args []string) error {
			cli := client.NewClient()
			return cli.StartClient(cmd, args)
		},
	}
	return cmd
}
