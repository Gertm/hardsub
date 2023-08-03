/*
		Copyright (C) 2023 Gert Meulyzer

	    This program is free software: you can redistribute it and/or modify
	    it under the terms of the GNU General Public License as published by
	    the Free Software Foundation, either version 3 of the License, or
	    (at your option) any later version.

	    This program is distributed in the hope that it will be useful,
	    but WITHOUT ANY WARRANTY; without even the implied warranty of
	    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	    GNU General Public License for more details.

	    You should have received a copy of the GNU General Public License
	    along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/
package main

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/bramvdbogaerde/go-scp"
	"github.com/bramvdbogaerde/go-scp/auth"
	"golang.org/x/crypto/ssh"
)

// remotePath needs to include the filename, I think. :-)
func CopyOverSCP(filename, host, username string, port int, remotePath string, privkeyPath string) error {
	if privkeyPath == "" {
		var err error
		privkeyPath, err = GetDefaultCertificate()
		if err != nil {
			return fmt.Errorf("don't have a valid certificate for scp'ing")
		}
	}
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

func GetDefaultCertificate() (string, error) {
	userHome, _ := os.UserHomeDir() // let's just assume we're on the right OS for the time being.
	defaultCert := path.Join(userHome, ".ssh", "id_rsa")
	if FileExists(defaultCert) {
		return defaultCert, nil
	}
	return "", fmt.Errorf("no default cert found")
}
