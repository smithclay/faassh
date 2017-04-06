package main

import (
	"flag"
	"net"

	"github.com/smithclay/tiny-ssh/server"
	"github.com/smithclay/tiny-ssh/tunnel"
	"golang.org/x/crypto/ssh"
)

var (
	sshdPort           = flag.String("port", "2200", "Port number for the ssh server to listen on")
	jumpHost           = flag.String("jh", "localhost", "Jump host")
	jumpHostPort       = flag.String("jh-port", "22", "Jump host SSH port number")
	jumpHostUser       = flag.String("jh-user", "ec2-user", "Jump host SSH user")
	jumpHostTunnelPort = flag.String("tunnel-port", "5001", "Jump host tunnel port")

	hostPrivateKey = flag.String("i", "id_rsa", "Path to RSA host private key")
)

func hostKeyCallback(hostname string, remote net.Addr, key ssh.PublicKey) error {
	return nil
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

	remoteEndpoint := &tunnel.Endpoint{
		HostPort: net.JoinHostPort("127.0.0.1", *jumpHostTunnelPort),
	}

	// Only key authentication is supported at this point.
	sshTunnelConfig := &ssh.ClientConfig{
		User: serverEndpoint.User,
		Auth: []ssh.AuthMethod{
			tunnel.SSHAgent(*hostPrivateKey),
		},
		HostKeyCallback: hostKeyCallback,
	}

	t := &tunnel.SSHtunnel{
		Config: sshTunnelConfig,
		Local:  localEndpoint,
		Server: serverEndpoint,
		Remote: remoteEndpoint,
	}
	go t.Start()

	// Create SSH Server with Dumb Terminal
	s := &server.SecureServer{
		User:     "foo",
		Password: "bar",
		HostKey:  *hostPrivateKey,
		Port:     *sshdPort,
	}
	s.Start()
}
