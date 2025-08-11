package ipset

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// CreateSet creates an ipset with the given name and type.
// setType can be for example "hash:ip" or "hash:net".
// If the set already exists, the command succeeds due to the
// use of the -exist flag.
func CreateSet(name, setType string) error {
	cmd := exec.Command("ipset", "create", name, setType, "-exist")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("ipset create failed: %v (%s)", err, string(output))
	}
	return nil
}

// AddEntry adds an entry (IP address or network) to the specified set.
// The -exist flag prevents an error if the entry is already present.
func AddEntry(setName, entry string) error {
	cmd := exec.Command("ipset", "add", setName, entry, "-exist")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("ipset add failed: %v (%s)", err, string(output))
	}
	return nil
}

// ListEntries returns all entries from the specified set.
func ListEntries(setName string) ([]string, error) {
	cmd := exec.Command("ipset", "list", setName)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("ipset list failed: %w", err)
	}

	var entries []string
	scanner := bufio.NewScanner(&buf)
	start := false
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Members:") {
			start = true
			continue
		}
		if start {
			line = strings.TrimSpace(line)
			if line != "" {
				entries = append(entries, line)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return entries, nil
}

// DeleteEntry removes an entry from the specified set.
func DeleteEntry(setName, entry string) error {
	cmd := exec.Command("ipset", "del", setName, entry)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("ipset del failed: %v (%s)", err, string(output))
	}
	return nil
}
