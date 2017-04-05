package main

// Originally from: https://gist.github.com/jpillora/b480fde82bff51a06238
// Port forwarding from: https://gist.github.com/ir4y/11146415
// and https://gist.github.com/sohlich/d8fb946f30a38d5a19f960a03ec1d740

import (
	"flag"
	"net"

	"github.com/smithclay/tiny-ssh/server"
	"github.com/smithclay/tiny-ssh/tunnel"
	"golang.org/x/crypto/ssh"
)

var (
	sshPort        = flag.String("port", "2200", "Port number to listen on")
	hostPrivateKey = flag.String("i", "id_rsa", "Path to RSA host private key")
)

func hostKeyCallback(hostname string, remote net.Addr, key ssh.PublicKey) error {
	return nil
}

func main() {
	flag.Parse()
	// Create SSH Tunnel
	localEndpoint := &tunnel.Endpoint{
		Host: "127.0.0.1",
		Port: "3000",
	}

	serverEndpoint := &tunnel.Endpoint{
		Host: "52.42.5.62",
		User: "ec2-user",
		Port: "22",
	}

	remoteEndpoint := &tunnel.Endpoint{
		Host: "127.0.0.1",
		Port: "5001",
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

	s := &server.SecureServer{
		User:     "foo",
		Password: "bar",
		HostKey:  *hostPrivateKey,
		Port:     *sshPort,
	}
	s.Start()
}
