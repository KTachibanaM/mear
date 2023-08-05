package ssh

import (
	"bufio"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

var SshProbingMaxAttempts = 30
var SshProbingInterval = 10 * time.Second

type TerminateSsh chan *SshStatus

func stream_scanner(scanner *bufio.Scanner, terminate_ssh TerminateSsh) {
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			agent_failure := print_log_and_parse_agent_failure(line)
			if agent_failure != nil {
				terminate_ssh <- NewSshAgentFailure(agent_failure)
			}
		}
	}
}

func SshExec(hostname, username string, private_key []byte, command string, timeout time.Duration) *SshStatus {
	// Parse private key
	signer, err := ssh.ParsePrivateKey(private_key)
	if err != nil {
		return NewSshError(fmt.Errorf("failed to parse private key: %v", err))
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
	var conn *ssh.Client
	for i := 0; i < SshProbingMaxAttempts; i++ {
		conn, err = ssh.Dial("tcp", fmt.Sprintf("%v:22", hostname), config)
		if err == nil {
			break
		} else {
			log.Warnf("failed to connect to server: %v, waiting...", err)
			time.Sleep(SshProbingInterval)
		}
	}
	if conn == nil {
		return NewSshError(fmt.Errorf("failed to connect to server: %v", err))
	}
	defer conn.Close()

	// Create new SSH session
	session, err := conn.NewSession()
	if err != nil {
		return NewSshError(fmt.Errorf("failed to create ssh session: %v", err))
	}
	defer session.Close()

	// Set up pipes for stdout and stderr
	stdout, err := session.StdoutPipe()
	if err != nil {
		return NewSshError(fmt.Errorf("failed to create stdout pipe: %v", err))
	}
	stdout_scanner := bufio.NewScanner(stdout)
	stderr, err := session.StderrPipe()
	if err != nil {
		return NewSshError(fmt.Errorf("failed to create stderr pipe: %v", err))
	}
	stderr_scanner := bufio.NewScanner(stderr)

	// Start remote command
	err = session.Start(command)
	if err != nil {
		return NewSshError(fmt.Errorf("failed to start command: %v", err))
	}

	terminate_ssh := make(chan *SshStatus)

	// Continuously send the command's output over the channel
	go stream_scanner(stdout_scanner, terminate_ssh)
	go stream_scanner(stderr_scanner, terminate_ssh)

	// Wait for command to finish
	go func() {
		err = session.Wait()
		if err != nil {
			terminate_ssh <- NewSshError(err)
		} else {
			terminate_ssh <- NewSshSuccess()
		}
	}()

	return <-terminate_ssh
}
