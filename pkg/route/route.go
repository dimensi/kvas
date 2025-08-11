package route

import (
	"fmt"

	"github.com/coreos/go-iptables/iptables"
)

const (
	table = "nat"
	chain = "PREROUTING"
)

// AddTunnelRule adds iptables rule to match traffic destined for entries in the
// given ipset and coming from specified interface. The rule simply accepts this
// traffic. The interface is typically a tunnel interface like "tun0".
func AddTunnelRule(iface, setName string) error {
	ipt, err := iptables.New()
	if err != nil {
		return fmt.Errorf("init iptables: %w", err)
	}

	ruleSpec := []string{"-i", iface, "-m", "set", "--match-set", setName, "dst", "-j", "ACCEPT"}
	exists, err := ipt.Exists(table, chain, ruleSpec...)
	if err != nil {
		return fmt.Errorf("check rule: %w", err)
	}
	if !exists {
		if err := ipt.Append(table, chain, ruleSpec...); err != nil {
			return fmt.Errorf("append rule: %w", err)
		}
	}
	return nil
}

// DeleteTunnelRule removes previously added rule from iptables.
func DeleteTunnelRule(iface, setName string) error {
	ipt, err := iptables.New()
	if err != nil {
		return fmt.Errorf("init iptables: %w", err)
	}
	ruleSpec := []string{"-i", iface, "-m", "set", "--match-set", setName, "dst", "-j", "ACCEPT"}
	exists, err := ipt.Exists(table, chain, ruleSpec...)
	if err != nil {
		return fmt.Errorf("check rule: %w", err)
	}
	if exists {
		if err := ipt.Delete(table, chain, ruleSpec...); err != nil {
			return fmt.Errorf("delete rule: %w", err)
		}
	}
	return nil
}
