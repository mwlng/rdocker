package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"os/user"
	"time"

	"rdocker/pkg/sshtunnel"

	"golang.org/x/crypto/ssh"
)

const (
	listenPort = 2308
)

var (
	username string
	keyFile  string
	tgtHost  string

	localEndpoint  sshtunnel.Endpoint
	serverEndpoint sshtunnel.Endpoint
	remoteEndpoint sshtunnel.Endpoint
	sshConfig      ssh.ClientConfig
	tunnel         sshtunnel.SSHTunnel
)

func defineFlags() {
	flag.StringVar(&username, "u", "", "ssh user")
	flag.StringVar(&username, "user", "", "ssh user")
	flag.StringVar(&keyFile, "i", "", "ssh key file")
	flag.StringVar(&keyFile, "key", "", "ssh key file")
	flag.StringVar(&tgtHost, "H", "", "Target hostname")
	flag.StringVar(&tgtHost, "host", "", "Target hostname")
}

func init() {
	defineFlags()

	localEndpoint = sshtunnel.Endpoint{
		Proto: "tcp",
		Addr: sshtunnel.IPAddr{
			Host: "localhost",
			Port: listenPort,
		},
	}

	remoteEndpoint = sshtunnel.Endpoint{
		Proto: "unix",
		Addr: sshtunnel.UnixAddr{
			SockPath: "/var/run/docker.sock",
		},
	}

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	keyFile = usr.HomeDir + "/.ssh/id_rsa"
}

func validateFlags() error {
	if tgtHost == "" {
		return errors.New("Missing -H/--host flag")
	}
	if username == "" {
		return errors.New("Missing -u/--user flag")
	}
	return nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func main() {
	flag.Parse()
	if err := validateFlags(); err != nil {
		fmt.Printf("Error: %s\n", err)
		flag.Usage()
		os.Exit(1)
	}

	var args []string
	for i, arg := range os.Args {
		if arg == "--" {
			args = os.Args[i+1:]
			break
		}
	}

	if len(args) == 0 {
		fmt.Println("Error: Missing docker command")
		os.Exit(1)
	}

	if !fileExists(keyFile) {
		fmt.Printf("Error: %s file does not exist\n", keyFile)
		os.Exit(1)
	}

	key, err := ioutil.ReadFile(keyFile)
	if err != nil {
		log.Fatal(err)
	}

	prikey, err := ssh.ParsePrivateKey([]byte(key))
	if err != nil {
		log.Fatal(err)
	}

	sshConfig = ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{ssh.PublicKeys(prikey)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Timeout: time.Second * 15,
	}

	serverEndpoint = sshtunnel.Endpoint{
		Proto: "tcp",
		Addr: sshtunnel.IPAddr{
			Host: tgtHost,
			Port: 22,
		},
	}

	tunnel := &sshtunnel.SSHTunnel{
		Local:  &localEndpoint,
		Server: &serverEndpoint,
		Remote: &remoteEndpoint,
		Config: &sshConfig,
	}

	go tunnel.Start()

	time.Sleep(time.Second * 2)

	args = append([]string{"-H", fmt.Sprintf("localhost:%d", listenPort)}, args...)
	cmd := exec.Command("docker", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
