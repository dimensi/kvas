package dns

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	// ConfigFile is the dnsmasq configuration file managed by kvasx
	ConfigFile = "/etc/dnsmasq.d/kvasx.conf"
)

// GenerateHosts creates a hosts file from the provided domain to IP mapping.
// Each entry is written as "IP domain" which is the format expected by dnsmasq.
func GenerateHosts(entries map[string]string, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for host, ip := range entries {
		if _, err := fmt.Fprintf(w, "%s %s\n", ip, host); err != nil {
			return err
		}
	}
	return w.Flush()
}

// GenerateIPSet writes dnsmasq ipset configuration file. Each domain is mapped
// to the provided ipset name.
func GenerateIPSet(domains []string, setName, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for _, d := range domains {
		if _, err := fmt.Fprintf(w, "ipset=/%s/%s\n", d, setName); err != nil {
			return err
		}
	}
	return w.Flush()
}

// SetServer updates upstream DNS server in ConfigFile.
func SetServer(server string) error {
	return updateConfigLine("server=", fmt.Sprintf("server=%s", server))
}

// SetPort updates listening port in ConfigFile.
func SetPort(port int) error {
	return updateConfigLine("port=", fmt.Sprintf("port=%d", port))
}

// updateConfigLine ensures the configuration file contains exactly one line
// starting with prefix. Existing lines with the same prefix are removed before
// the new line is appended.
func updateConfigLine(prefix, line string) error {
	var lines []string
	if data, err := os.ReadFile(ConfigFile); err == nil {
		for _, l := range strings.Split(string(data), "\n") {
			if !strings.HasPrefix(l, prefix) && l != "" {
				lines = append(lines, l)
			}
		}
	}
	lines = append(lines, line)
	return os.WriteFile(ConfigFile, []byte(strings.Join(lines, "\n")+"\n"), 0644)
}

// removeConfigLine removes all lines starting with prefix from ConfigFile.
func removeConfigLine(prefix string) error {
	data, err := os.ReadFile(ConfigFile)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	var lines []string
	if err == nil {
		for _, l := range strings.Split(string(data), "\n") {
			if !strings.HasPrefix(l, prefix) && l != "" {
				lines = append(lines, l)
			}
		}
	}
	if len(lines) == 0 {
		return os.Remove(ConfigFile)
	}
	return os.WriteFile(ConfigFile, []byte(strings.Join(lines, "\n")+"\n"), 0644)
}
