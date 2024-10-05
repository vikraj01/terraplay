package utils

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/ssh"
)

type SSHConfig struct {
	Host       string
	Port       string
	User       string
	PrivateKey []byte
}

func ConnectToEC2ViaSSH(config SSHConfig) (*ssh.Client, error) {
	signer, err := ssh.ParsePrivateKey(config.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("error parsing SSH private key: %v", err)
	}

	sshConfig := &ssh.ClientConfig{
		User:            config.User,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", config.Host, config.Port), sshConfig)
	if err != nil {
		return nil, fmt.Errorf("error dialing SSH: %v", err)
	}

	return client, nil
}

func RunCommandOnEC2(client *ssh.Client, command string) error {
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("error creating SSH session: %v", err)
	}
	defer session.Close()

	var stdout, stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	if err := session.Run(command); err != nil {
		return fmt.Errorf("error running command: %v (stderr: %s)", err, stderr.String())
	}

	log.Printf("Command output: %s", stdout.String())
	return nil
}

func ConnectToEC2ViaSSHWithRetry(sshConfig SSHConfig) (*ssh.Client, error) {
	signer, err := ssh.ParsePrivateKey(sshConfig.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("error parsing SSH private key: %v", err)
	}

	sshClientConfig := &ssh.ClientConfig{
		User: sshConfig.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), 
	}

	for retries := 0; retries < 12; retries++ {
		client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", sshConfig.Host, sshConfig.Port), sshClientConfig)
		if err == nil {
			log.Printf("Successfully connected to EC2 via SSH: %s", sshConfig.Host)
			return client, nil
		}

		log.Printf("SSH connection to EC2 failed (retry %d/12): %v", retries+1, err)
		time.Sleep(10 * time.Second)
	}

	return nil, fmt.Errorf("failed to connect to EC2 via SSH after multiple attempts")
}