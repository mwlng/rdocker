package sshtunnel

import (
	"fmt"
	"io"
	"net"
	"os"
	"sync"

	"golang.org/x/crypto/ssh"
)

type SSHTunnel struct {
	Local  *Endpoint
	Server *Endpoint
	Remote *Endpoint

	Config *ssh.ClientConfig
}

func (tunnel *SSHTunnel) Start() error {
	localAddr, err := tunnel.Local.ToString()
	if err != nil {
		return err
	}
	listener, err := net.Listen("tcp", localAddr)
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go tunnel.forward(conn)
	}
}

func (tunnel *SSHTunnel) forward(localConn net.Conn) {
	serverAddr, err := tunnel.Server.ToString()
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	remoteAddr, err := tunnel.Remote.ToString()
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	serverConn, err := ssh.Dial("tcp", serverAddr, tunnel.Config)
	if err != nil {
		fmt.Printf("Server dial error: %s\n", err)
		os.Exit(1)
	}

	remoteConn, err := serverConn.Dial("unix", remoteAddr)
	if err != nil {
		fmt.Printf("Remote dial error: %s\n", err)
		os.Exit(1)
	}
	defer localConn.Close()
	defer remoteConn.Close()

	copyConn := func(writer, reader net.Conn, wg sync.WaitGroup) {
		defer wg.Done()
		_, err := io.Copy(writer, reader)
		if err != nil {
			if err == io.EOF {
				reader.Close()
				return
			} else {
				fmt.Printf("io.Copy error: %s", err)
			}
		}
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go copyConn(localConn, remoteConn, wg)

	wg.Add(1)
	go copyConn(remoteConn, localConn, wg)

	wg.Wait()
}
