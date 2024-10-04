package utils

import (
	"bytes"
	"fmt"
	"log"
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