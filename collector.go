package main

import (
	"bytes"
	"fmt"
	"log"
	"regexp"
	"strings"

	"golang.org/x/crypto/ssh"
)

func main() {

	server := "134.119.179.22"
	port := 8822
	user := "zabbix"
	password := "RevDau@123"

	authMethods := []ssh.AuthMethod{ssh.Password(password)}

	config := &ssh.ClientConfig{
		User: user,
		Auth: authMethods,

		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", server, port), config)
	if err != nil {
		log.Fatal("failed to dial:", err)
	}
	defer client.Close()

	// Create a new session
	session, err := client.NewSession()
	if err != nil {
		log.Fatal("failed to create session:", err)
	}
	defer session.Close()

	// Capture standard output and standard error
	var stdout, stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	command := "ifconfig"

	err = session.Run(command)
	if err != nil {
		log.Printf("remote command failed: %v\n", err)
		fmt.Println("stderr:", string(stderr.Bytes()))
	}

	interfaceRegex := regexp.MustCompile(`^([a-zA-Z0-9]+):`)
	ipv4Regex := regexp.MustCompile(`inet ([0-9.]+)`)
	ipv6Regex := regexp.MustCompile(`inet6 ([0-9a-fA-F:]+)`)

	// Split output into lines
	lines := strings.Split(string(stdout.Bytes()), "\n")

	//variables to store interface and addresses
	var interfaceName string
	var ipv4Addr, ipv6Addr string

	// Iterate over lines to find interfaces and their addresses
	for _, line := range lines {
		if matches := interfaceRegex.FindStringSubmatch(line); len(matches) > 1 {

			if interfaceName != "" {
				printInterfaceInfo(interfaceName, ipv4Addr, ipv6Addr)
			}

			// Update interface name
			interfaceName = matches[1]
			// Reset addresses
			ipv4Addr = ""
			ipv6Addr = ""
		} else if ipv4Matches := ipv4Regex.FindStringSubmatch(line); len(ipv4Matches) > 1 {
			ipv4Addr = ipv4Matches[1]
		} else if ipv6Matches := ipv6Regex.FindStringSubmatch(line); len(ipv6Matches) > 1 {
			ipv6Addr = ipv6Matches[1]
		}
	}

	// Print information for the last interface
	if interfaceName != "" {
		printInterfaceInfo(interfaceName, ipv4Addr, ipv6Addr)
	}
}

// function to print interface information
func printInterfaceInfo(interfaceName, ipv4Addr, ipv6Addr string) {
	fmt.Printf("Interface: %s\n", interfaceName)
	if ipv4Addr != "" {
		fmt.Printf("IPv4 Address: %s\n", ipv4Addr)
	}
	if ipv6Addr != "" {
		fmt.Printf("IPv6 Address: %s\n", ipv6Addr)
	}
	fmt.Println()
}
