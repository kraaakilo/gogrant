package main

import (
	"bufio"
	"strings"
	"testing"
)

func TestPromptNetworkIPLastOctet(t *testing.T) {
	// Simulate user input for last octet "50"
	input := strings.NewReader("50\n")
	reader := bufio.NewReader(input)

	ip := promptNetworkIP(reader, "Test IP", "192.168.56", "100")
	expected := "192.168.56.50"

	if ip != expected {
		t.Errorf("expected %s, got %s", expected, ip)
	}
}

func TestPromptNetworkIPFullIP(t *testing.T) {
	// Simulate user input for full IP
	input := strings.NewReader("10.0.0.42\n")
	reader := bufio.NewReader(input)

	ip := promptNetworkIP(reader, "Test IP", "192.168.56", "100")
	expected := "10.0.0.42"

	if ip != expected {
		t.Errorf("expected %s, got %s", expected, ip)
	}
}
