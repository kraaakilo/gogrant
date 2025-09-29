package main

import (
	"bufio"
	"embed"
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
	"text/template"
)

type Template struct {
	BoxName       string
	Hostname      string
	Memory        int
	CpuCount      int
	EnablePrivate bool
	EnablePublic  bool
	PrivateIP     string
	PublicIP      string
	PrivateBase   string
	PublicBase    string
}

//go:embed static/*
var staticFS embed.FS

func main() {
	if err := generate(); err != nil {
		log.Fatal(err)
	}
}

func generate() error {
	content, err := staticFS.ReadFile("static/Vagrantfile")
	if err != nil {
		return fmt.Errorf("reading template: %w", err)
	}

	reader := bufio.NewReader(os.Stdin)

	base := Template{
		BoxName:       "ubuntu/jammy64",
		Hostname:      "gogrant-box",
		Memory:        1024,
		CpuCount:      2,
		EnablePrivate: true,
		EnablePublic:  false,
		PrivateBase:   "192.168.56",
		PublicBase:    "192.168.1",
	}

	// check cpu count
	if runtime.NumCPU() <= int(base.CpuCount) {
		return fmt.Errorf("not enough CPUs: %d required, %d available", base.CpuCount, runtime.NumCPU())
	}

	// interactive prompts
	base.BoxName = promptString(reader, "Box image name", base.BoxName)
	base.Hostname = promptString(reader, "Hostname", base.Hostname)
	base.Memory = promptInt(reader, "Memory (MB)", base.Memory)
	base.CpuCount = promptInt(reader, "CPU count", base.CpuCount)
	base.EnablePrivate = promptBool(reader, "Enable private network? (y/N)", base.EnablePrivate)
	if base.EnablePrivate {
		base.PrivateIP = promptNetworkIP(reader, "Private network IP", base.PrivateBase, "100")
	}

	base.EnablePublic = promptBool(reader, "Enable public network? (y/N)", base.EnablePublic)
	if base.EnablePublic {
		base.PublicIP = promptNetworkIP(reader, "Public network IP", base.PublicBase, "50")
	}

	// parse template
	tmpl, err := template.New("vagrant").Parse(string(content))
	if err != nil {
		return fmt.Errorf("parsing template: %w", err)
	}

	// write vagrantfile
	file, err := os.Create("Vagrantfile")
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	}
	defer file.Close()

	if err := tmpl.Execute(file, base); err != nil {
		return fmt.Errorf("executing template: %w", err)
	}

	fmt.Println("[OK] Vagrantfile successfully generated.")
	return nil
}

// helpers

func promptString(reader *bufio.Reader, label string, def string) string {
	fmt.Printf("%s (default: %s): ", label, def)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		return def
	}
	return input
}

func promptInt(reader *bufio.Reader, label string, def int) int {
	for {
		fmt.Printf("%s (default: %d): ", label, def)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "" {
			return def
		}
		val, err := strconv.Atoi(input)
		if err == nil {
			return val
		}
		fmt.Println("Invalid number, try again.")
	}
}

func promptBool(reader *bufio.Reader, label string, def bool) bool {
	fmt.Printf("%s (default: %v): ", label, def)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))
	if input == "" {
		return def
	}
	return input == "y" || input == "yes"
}

// prompt for ip: last octet or full ip
func promptNetworkIP(reader *bufio.Reader, label, base, defaultOctet string) string {
	for {
		fmt.Printf("%s (default: %s.%s): ", label, base, defaultOctet)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "" {
			return base + "." + defaultOctet
		}

		if strings.Count(input, ".") == 0 {
			if octet, err := strconv.Atoi(input); err == nil && octet >= 0 && octet <= 255 {
				return base + "." + input
			}
			fmt.Println("Invalid octet (0-255), try again.")
			continue
		}

		if ip := net.ParseIP(input); ip != nil {
			return input
		}

		fmt.Println("Invalid IP format, try again.")
	}
}
