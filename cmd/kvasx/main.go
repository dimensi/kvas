package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	cfgPath  string
	logLevel string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "kvasx",
		Short: "kvasx is a CLI tool",
	}

	rootCmd.PersistentFlags().StringVar(&cfgPath, "config", "", "Path to config file")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "Logging level")

	dnsCmd := &cobra.Command{
		Use:   "dns",
		Short: "DNS related commands",
	}

	dnsStatusCmd := &cobra.Command{
		Use:   "status",
		Short: "Show DNS status",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("DNS status: OK")
		},
	}

	dnsCmd.AddCommand(dnsStatusCmd)
	rootCmd.AddCommand(dnsCmd)

	rootCmd.AddCommand(&cobra.Command{
		Use:   "vpn",
		Short: "VPN related commands",
	})
	rootCmd.AddCommand(&cobra.Command{
		Use:   "adblock",
		Short: "Adblock related commands",
	})

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
