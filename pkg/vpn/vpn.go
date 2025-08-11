package vpn

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	// ConfigFile is the path to the generated xray configuration.
	ConfigFile = "/opt/etc/xray/kvas.json"
	// binary name for xray/v2ray daemon
	bin = "xray"
)

// Server describes remote server parameters for VLESS/VMess protocols.
type Server struct {
	Address string
	Port    int
	UUID    string
	Domain  string // SNI domain
}

// GenerateConfig builds a minimal VLESS client configuration for xray/v2ray.
func GenerateConfig(s Server) ([]byte, error) {
	conf := map[string]any{
		"log": map[string]any{
			"loglevel": "warning",
		},
		"inbounds": []map[string]any{
			{
				"port":     10808,
				"listen":   "127.0.0.1",
				"protocol": "socks",
				"settings": map[string]any{},
			},
		},
		"outbounds": []map[string]any{
			{
				"protocol": "vless",
				"settings": map[string]any{
					"vnext": []map[string]any{
						{
							"address": s.Address,
							"port":    s.Port,
							"users": []map[string]any{
								{
									"id":         s.UUID,
									"encryption": "none",
								},
							},
						},
					},
				},
				"streamSettings": map[string]any{
					"network":  "tcp",
					"security": "tls",
					"tlsSettings": map[string]any{
						"serverName": s.Domain,
					},
				},
			},
		},
	}
	return json.MarshalIndent(conf, "", "  ")
}

// WriteConfig generates and writes xray configuration to the provided path.
func WriteConfig(s Server, path string) error {
	data, err := GenerateConfig(s)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// Start launches the xray daemon with the given configuration path.
func Start(configPath string) error {
	cmd := exec.Command(bin, "run", "-config", configPath)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start %s: %w", bin, err)
	}
	return cmd.Process.Release()
}

// Stop terminates the running xray daemon.
func Stop() error {
	cmd := exec.Command("pkill", bin)
	if output, err := cmd.CombinedOutput(); err != nil {
		if len(output) > 0 {
			return fmt.Errorf("pkill %s: %v (%s)", bin, err, strings.TrimSpace(string(output)))
		}
		return fmt.Errorf("pkill %s: %w", bin, err)
	}
	return nil
}

// IsRunning returns true if the xray daemon is running.
func IsRunning() bool {
	cmd := exec.Command("pidof", bin)
	return cmd.Run() == nil
}

// CheckDomain performs HTTP GET request and measures duration.
func CheckDomain(domain string) (time.Duration, error) {
	if !strings.HasPrefix(domain, "http://") && !strings.HasPrefix(domain, "https://") {
		domain = "https://" + domain
	}
	start := time.Now()
	resp, err := http.Get(domain)
	if err != nil {
		return 0, err
	}
	resp.Body.Close()
	return time.Since(start), nil
}

// Status reports whether daemon is running and measures availability of domain.
func Status(domain string) (bool, time.Duration, error) {
	running := IsRunning()
	if !running {
		return false, 0, nil
	}
	d, err := CheckDomain(domain)
	return true, d, err
}

// PromptServer asks user for server parameters via stdin.
func PromptServer(r *bufio.Reader) (Server, error) {
	var s Server
	fmt.Print("Server address: ")
	addr, err := r.ReadString('\n')
	if err != nil {
		return s, err
	}
	fmt.Print("Server port: ")
	portStr, err := r.ReadString('\n')
	if err != nil {
		return s, err
	}
	fmt.Print("User UUID: ")
	uuid, err := r.ReadString('\n')
	if err != nil {
		return s, err
	}
	fmt.Print("SNI domain: ")
	dom, err := r.ReadString('\n')
	if err != nil {
		return s, err
	}
	addr = strings.TrimSpace(addr)
	portStr = strings.TrimSpace(portStr)
	uuid = strings.TrimSpace(uuid)
	dom = strings.TrimSpace(dom)
	var port int
	fmt.Sscanf(portStr, "%d", &port)
	s = Server{Address: addr, Port: port, UUID: uuid, Domain: dom}
	return s, nil
}
