package ssh

import (
	"bufio"
	"fmt"
	"time"

	"golang.org/x/crypto/ssh"
)

func stream_scanner(scanner *bufio.Scanner) {
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			print_log(line)
		}
	}
}

func SshExec(hostname, username string, private_key []byte, command string, timeout time.Duration) error {
	// Parse private key
	signer, err := ssh.ParsePrivateKey(private_key)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %v", err)
	}

	// Set up SSH client configuration with private key
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         timeout,
	}

	// Connect to remote server
	conn, err := ssh.Dial("tcp", fmt.Sprintf("%v:22", hostname), config)
	if err != nil {
		return fmt.Errorf("failed to connect to server: %v", err)
	}
	defer conn.Close()

	// Create new SSH session
	session, err := conn.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create ssh session: %v", err)
	}
	defer session.Close()

	// Set up pipes for stdout and stderr
	stdout, err := session.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %v", err)
	}
	stderr, err := session.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %v", err)
	}

	// Start remote command
	err = session.Start(command)
	if err != nil {
		return fmt.Errorf("failed to start command: %v", err)
	}

	// Continuously send the command's output over the channel
	go stream_scanner(bufio.NewScanner(stdout))
	go stream_scanner(bufio.NewScanner(stderr))

	// Wait for command to finish
	err = session.Wait()
	if err != nil {
		return fmt.Errorf("command failed: %v", err)
	}

	return nil
}
