package main

import (
	"bufio"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/smithclay/faassh/server"
	"github.com/smithclay/faassh/tunnel"
	"golang.org/x/crypto/ssh"
)

var (
	sshdPort           = flag.String("port", "2200", "Port number for the ssh server to listen on")
	jumpHost           = flag.String("jh", "localhost", "Jump host")
	jumpHostPort       = flag.String("jh-port", "22", "Jump host SSH port number")
	jumpHostUser       = flag.String("jh-user", "ec2-user", "Jump host SSH user")
	jumpHostTunnelPort = flag.String("tunnel-port", "0", "Jump host tunnel port")

	hostPrivateKey = flag.String("i", "id_rsa", "Path to RSA host private key")
)

func hostKeyCallback(hostname string, remote net.Addr, key ssh.PublicKey) error {
	return nil
}

func readStdin(s *bufio.Scanner) {
	for s.Scan() {
		log.Println("line", s.Text())
	}
}

func main() {
	flag.Parse()

	// Create SSH Tunnel
	localEndpoint := &tunnel.Endpoint{
		HostPort: net.JoinHostPort("127.0.0.1", *sshdPort),
	}

	serverEndpoint := &tunnel.Endpoint{
		HostPort: net.JoinHostPort(*jumpHost, *jumpHostPort),
		User:     *jumpHostUser,
	}

	// With the '0' default, an open port on the host will be chosen automatically.
	remoteEndpoint := &tunnel.Endpoint{
		HostPort: net.JoinHostPort("127.0.0.1", *jumpHostTunnelPort),
	}

	// Only key authentication is supported at this point.
	sshTunnelConfig := &ssh.ClientConfig{
		User: serverEndpoint.User,
		Auth: []ssh.AuthMethod{
			tunnel.SSHAgent(*hostPrivateKey),
		},
		Timeout:         time.Second * 10,
		HostKeyCallback: hostKeyCallback,
	}
	// Create SSH Server with Dumb Terminal
	s := &server.SecureServer{
		User:     "foo",
		Password: "bar",
		HostKey:  *hostPrivateKey,
		Port:     *sshdPort,
	}

	t := &tunnel.SSHtunnel{
		Config: sshTunnelConfig,
		Local:  localEndpoint,
		Server: serverEndpoint,
		Remote: remoteEndpoint,
	}

	scanner := bufio.NewScanner(os.Stdin)
	go readStdin(scanner)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)
	go func() {
		sig := <-sigs
		log.Printf("%v: Attempting to stop server and close tunnel...", sig)

		sErr := s.Stop()
		if sErr != nil {
			log.Printf("Could not stop ssh server: %v", sErr)
		}
		tErr := t.Stop()
		if tErr != nil {
			log.Printf("Could not stop tunnel: %v", tErr)
		}

		if tErr != nil || sErr != nil {
			os.Exit(1)
		}
		os.Exit(0)
	}()

	go t.Start()
	s.Start()
}
