package tunnel

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"

	"golang.org/x/crypto/ssh"
)

type Endpoint struct {
	Host string
	User string
	Port string
}

func (endpoint *Endpoint) String() string {
	return fmt.Sprintf("%s:%s", endpoint.Host, endpoint.Port)
}

type SSHtunnel struct {
	Local  *Endpoint
	Server *Endpoint
	Remote *Endpoint

	Config *ssh.ClientConfig
}

func (t *SSHtunnel) Start() error {
	log.Printf("Creating tunnel to %v...", t.Server.String())
	conn, err := ssh.Dial("tcp", t.Server.String(), t.Config)
	if err != nil {
		log.Fatalf("unable to connect to remote server: %v", err)
	}
	defer conn.Close()

	log.Printf("Registering tcp forward on %v", t.Remote.String())
	remoteListener, err := conn.Listen("tcp", t.Remote.String())
	if err != nil {
		log.Fatalf("unable to register tcp forward: %v", err)
	}
	defer remoteListener.Close()

	for {
		r, err := remoteListener.Accept()
		if err != nil {
			log.Fatalf("listen.Accept failed: %v", err)
		}
		go t.forward(r)
	}
}

func (t *SSHtunnel) forward(remoteConn net.Conn) {
	log.Printf("Registering local tcp forward on %v", t.Local.String())
	localConn, err := net.Dial("tcp", t.Local.String())
	if err != nil {
		log.Fatalf("local: unable to register tcp forward: %v", err)
	}

	copyConn := func(writer, reader net.Conn) {
		_, err := io.Copy(writer, reader)
		if err != nil {
			fmt.Printf("io.Copy error: %s", err)
		}
	}

	go copyConn(localConn, remoteConn)
	go copyConn(remoteConn, localConn)
}

func SSHAgent(keyfile string) ssh.AuthMethod {
	key, err := ioutil.ReadFile(keyfile)
	if err != nil {
		log.Fatalf("unable to read private key: %v", err)
	}
	signer, err := ssh.ParsePrivateKey(key)

	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
	}
	return ssh.PublicKeys(signer)
}
