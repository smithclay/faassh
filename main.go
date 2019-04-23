package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/smithclay/faassh/server"
	"github.com/smithclay/faassh/tunnel"
	"golang.org/x/crypto/ssh"
)

var (
	tunnelCommand = flag.NewFlagSet("tunnel", flag.ExitOnError)
)

var (
	sshdPort           = flag.String("port", "2200", "Port number for ssh server (non-priviliged)")
	hostPrivateKey     = flag.String("i", "", "Path to RSA host private key")
	jumpHost           = tunnelCommand.String("jh", "localhost", "Jump host")
	jumpHostPort       = tunnelCommand.String("jh-port", "22", "Jump host SSH port number")
	jumpHostUser       = tunnelCommand.String("jh-user", "ec2-user", "Jump host SSH user")
	jumpHostTunnelPort = tunnelCommand.String("tunnel-port", "0", "Jump host tunnel port")
)

// Only key authentication is supported at this point.
// This will accept connections from any remote host.
func hostKeyCallback(hostname string, remote net.Addr, key ssh.PublicKey) error {
	return nil
}

func createTunnel(localPort string, jumpHost string, jumpHostPort string, jumpHostUser string, jumpHostTunnelPort string) *tunnel.SSHtunnel {
	// Create SSH Tunnel
	// Example: 127.0.0.1:2200
	localEndpoint := &tunnel.Endpoint{
		HostPort: net.JoinHostPort("127.0.0.1", localPort),
	}
	// Jump Host Endpoint
	// Example: 0.tcp.ngrok.io:15303
	jumpEndpoint := &tunnel.Endpoint{
		HostPort: net.JoinHostPort(jumpHost, jumpHostPort),
		User:     jumpHostUser,
	}

	// With the '0' default, an open port on the host will be chosen automatically.
	// This is the endpoint the client (i.e. dev laptop) actually connects to.
	// Example: 127.0.0.1:5001
	// Then, `ssh -p 5001 foo@127.0.0.1` to connect to the function.
	remoteEndpoint := &tunnel.Endpoint{
		HostPort: net.JoinHostPort("127.0.0.1", jumpHostTunnelPort),
	}

	sshTunnelConfig := &ssh.ClientConfig{
		User: jumpEndpoint.User,
		Auth: []ssh.AuthMethod{
			tunnel.SSHAgent(*hostPrivateKey),
		},
		Timeout:         time.Second * 10,
		HostKeyCallback: hostKeyCallback,
	}

	return &tunnel.SSHtunnel{
		Config: sshTunnelConfig,
		Local:  localEndpoint,
		Server: jumpEndpoint,
		Remote: remoteEndpoint,
	}
}

func main() {
	flag.Parse()

	if *hostPrivateKey == "" {
		fmt.Println("Error: Please supply the host private key using the -i option.")
		os.Exit(1)
	}

	// Create SSH Server with Dumb Terminal
	s := &server.SecureServer{
		User:     "foo",
		Password: "bar",
		HostKey:  *hostPrivateKey,
		Port:     *sshdPort,
	}

	if tunnelCommand.Parsed() {
		t := createTunnel(*sshdPort, *jumpHost, *jumpHostPort, *jumpHostUser, *jumpHostTunnelPort)
		go t.Start()
	}

	s.Start()
}
