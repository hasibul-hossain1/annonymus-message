// anon - send anonymous messages to teammates on Telegram.
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type config struct {
	Server string `json:"server"`
	Key    string `json:"key"`
}

const usage = `anon - send anonymous messages to your team on Telegram

Usage:
  anon config <server-url> <team-key>   save server URL and team key
  anon send @username <message...>      send an anonymous message
  echo "msg" | anon send @username      message can also come from stdin
  anon members                          list registered teammates
  anon completion <zsh|bash|powershell> print shell completion script
`

var client = &http.Client{Timeout: 15 * time.Second}

func main() {
	if len(os.Args) < 2 {
		fmt.Print(usage)
		os.Exit(1)
	}

	var err error
	switch os.Args[1] {
	case "config":
		err = cmdConfig(os.Args[2:])
	case "send":
		err = cmdSend(os.Args[2:])
	case "members":
		err = cmdMembers()
	case "completion":
		err = cmdCompletion(os.Args[2:])
	case "__members":
		completeMembers()
	case "help", "-h", "--help":
		fmt.Print(usage)
	default:
		err = fmt.Errorf("unknown command %q\n\n%s", os.Args[1], usage)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func configPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "anon", "config.json"), nil
}

func loadConfig() (config, error) {
	var c config
	path, err := configPath()
	if err != nil {
		return c, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return c, errors.New("not configured, run: anon config <server-url> <team-key>")
	}
	err = json.Unmarshal(data, &c)
	return c, err
}

func cmdConfig(args []string) error {
	if len(args) != 2 {
		return errors.New("usage: anon config <server-url> <team-key>")
	}
	path, err := configPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	data, _ := json.MarshalIndent(config{Server: strings.TrimRight(args[0], "/"), Key: args[1]}, "", "  ")
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return err
	}
	fmt.Println("saved to", path)
	return nil
}

func cmdSend(args []string) error {
	if len(args) < 1 {
		return errors.New("usage: anon send @username <message...>")
	}
	to := args[0]
	message := strings.Join(args[1:], " ")

	if message == "" {
		if stat, _ := os.Stdin.Stat(); stat.Mode()&os.ModeCharDevice == 0 {
			data, err := io.ReadAll(os.Stdin)
			if err != nil {
				return err
			}
			message = string(data)
		}
	}
	if strings.TrimSpace(message) == "" {
		return errors.New("message is empty")
	}

	var res struct{}
	if err := call("POST", "/send", map[string]string{"to": to, "message": message}, &res); err != nil {
		return err
	}
	fmt.Println("✓ sent anonymously to", "@"+strings.TrimPrefix(to, "@"))
	return nil
}

func cmdMembers() error {
	members, err := fetchMembers()
	if err != nil {
		return err
	}
	if len(members) == 0 {
		fmt.Println("no one has registered yet (send /start to the bot)")
	}
	for _, m := range members {
		fmt.Println("@" + m)
	}
	return nil
}

// fetchMembers gets the member list from the server and refreshes the
// completion cache.
func fetchMembers() ([]string, error) {
	var res struct {
		Members []string `json:"members"`
	}
	if err := call("GET", "/members", nil, &res); err != nil {
		return nil, err
	}
	writeMembersCache(res.Members)
	return res.Members, nil
}

func call(method, path string, body any, out any) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	var reader io.Reader
	if body != nil {
		data, _ := json.Marshal(body)
		reader = bytes.NewReader(data)
	}
	req, err := http.NewRequest(method, cfg.Server+path, reader)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+cfg.Key)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		var e struct {
			Error string `json:"error"`
		}
		if json.Unmarshal(data, &e) == nil && e.Error != "" {
			return errors.New(e.Error)
		}
		return fmt.Errorf("server returned %s", resp.Status)
	}
	return json.Unmarshal(data, out)
}
