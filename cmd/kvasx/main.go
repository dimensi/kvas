package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"kvasx/pkg/dns"
	"kvasx/pkg/ipset"
	"kvasx/pkg/route"
	"kvasx/pkg/vpn"
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

	dnsSetCmd := &cobra.Command{
		Use:   "set <server>",
		Short: "Set upstream DNS server",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return dns.SetServer(args[0])
		},
	}

	dnsCmd.AddCommand(dnsStatusCmd)
	dnsCmd.AddCommand(dnsSetCmd)
	rootCmd.AddCommand(dnsCmd)

	vpnCmd := &cobra.Command{
		Use:   "vpn",
		Short: "VPN related commands",
	}

	vpnScanCmd := &cobra.Command{
		Use:   "scan",
		Short: "Scan VPN tunnel and configure iptables",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Adding iptables rule for VPN tunnel")
			// default interface and set name
			return route.AddTunnelRule("tun0", "kvas_vpn")
		},
	}

	vpnNewCmd := &cobra.Command{
		Use:   "new",
		Short: "Create new VPN configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			reader := bufio.NewReader(os.Stdin)
			s, err := vpn.PromptServer(reader)
			if err != nil {
				return err
			}
			if err := vpn.WriteConfig(s, vpn.ConfigFile); err != nil {
				return err
			}
			fmt.Printf("Config written to %s\n", vpn.ConfigFile)
			return nil
		},
	}

	var checkDomain string
	vpnStatusCmd := &cobra.Command{
		Use:   "status",
		Short: "Show VPN status and domain availability",
		RunE: func(cmd *cobra.Command, args []string) error {
			running, dur, err := vpn.Status(checkDomain)
			if err != nil {
				return err
			}
			if !running {
				fmt.Println("VPN daemon: stopped")
				return nil
			}
			fmt.Println("VPN daemon: running")
			fmt.Printf("Domain %s reachable in %v\n", checkDomain, dur)
			return nil
		},
	}
	vpnStatusCmd.Flags().StringVar(&checkDomain, "domain", "example.com", "Domain to test through VPN")

	vpnCmd.AddCommand(vpnScanCmd)
	vpnCmd.AddCommand(vpnNewCmd)
	vpnCmd.AddCommand(vpnStatusCmd)
	rootCmd.AddCommand(vpnCmd)

	ipsetCmd := &cobra.Command{
		Use:   "ipset",
		Short: "IPSet related commands",
	}
	ipsetAddCmd := &cobra.Command{
		Use:   "add <domain>",
		Short: "Add domain to ipset",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			domain := args[0]
			if err := ipset.CreateSet("kvas_vpn", "hash:ip"); err != nil {
				return err
			}
			fmt.Printf("Adding %s to ipset\n", domain)
			return ipset.AddEntry("kvas_vpn", domain)
		},
	}
	ipsetCmd.AddCommand(ipsetAddCmd)
	rootCmd.AddCommand(ipsetCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
