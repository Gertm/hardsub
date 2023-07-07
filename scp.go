package main

import (
	"context"
	"fmt"
	"os"

	"github.com/bramvdbogaerde/go-scp"
	"github.com/bramvdbogaerde/go-scp/auth"
	"golang.org/x/crypto/ssh"
)

// remotePath needs to include the filename, I think. :-)
func CopyOverSCP(filename, host, username string, port int, remotePath string, privkeyPath string) error {
	clientConfig, err := auth.PrivateKey(username, privkeyPath, ssh.InsecureIgnoreHostKey())
	if err != nil {
		return fmt.Errorf("cannot set up auth: %w", err)
	}
	client := scp.NewClient(fmt.Sprintf("%s:%d", host, port), &clientConfig)

	err = client.Connect()
	if err != nil {
		return fmt.Errorf("couldn't establish a connection to the remote server: %w", err)
	}

	f, _ := os.Open(filename)
	defer client.Close()
	defer f.Close()

	return client.CopyFromFile(context.Background(), *f, remotePath, "0655")
}
