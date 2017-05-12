package server

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os/exec"
	"strings"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

type SecureServer struct {
	User     string
	Password string
	HostKey  string
	Port     string
}

func (s *SecureServer) Stop() error {
	// TODO: Close all connections
	return nil
}

func (s *SecureServer) Start() error {
	// In the latest version of crypto/ssh (after Go 1.3), the SSH server type has been removed
	// in favour of an SSH connection type. A ssh.ServerConn is created by passing an existing
	// net.Conn and a ssh.ServerConfig to ssh.NewServerConn, in effect, upgrading the net.Conn
	// into an ssh.ServerConn

	config := &ssh.ServerConfig{
		//Define a function to run when a client attempts a password login
		PasswordCallback: func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
			// Should use constant-time compare (or better, salt+hash) in a production setting.
			if c.User() == s.User && string(pass) == s.Password {
				return nil, nil
			}
			return nil, fmt.Errorf("password rejected for %q", c.User())
		},
		// You may also explicitly allow anonymous client authentication, though anon bash
		// sessions may not be a wise idea
		// NoClientAuth: true,
	}

	// You can generate a keypair with 'ssh-keygen -t rsa'
	privateBytes, err := ioutil.ReadFile(s.HostKey)
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to load private key (%s): %s", s.HostKey, err))
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		log.Fatal("Failed to parse private key")
	}

	config.AddHostKey(private)

	// Once a ServerConfig has been configured, connections can be accepted.
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", s.Port))
	if err != nil {
		log.Fatalf("Failed to listen on %s (%s)", s.Port, err)
	}
	log.Print(fmt.Sprintf("Listening on %s...", s.Port))

	for {
		tcpConn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept incoming connection (%s)", err)
			continue
		}
		// Before use, a handshake must be performed on the incoming net.Conn.
		sshConn, chans, reqs, err := ssh.NewServerConn(tcpConn, config)
		if err != nil {
			log.Printf("Failed to handshake (%s)", err)
			continue
		}

		log.Printf("New SSH connection from %s (%s)", sshConn.RemoteAddr(), sshConn.ClientVersion())
		// Discard all global out-of-band Requests
		go ssh.DiscardRequests(reqs)

		// Accept all channels
		for newChannel := range chans {
			go s.handleChannel(newChannel)
		}
	}
}

func (s *SecureServer) handleChannel(newChannel ssh.NewChannel) {
	// Since we're handling a shell, we expect a
	// channel type of "session". The also describes
	// "x11", "direct-tcpip" and "forwarded-tcpip"
	// channel types.
	if t := newChannel.ChannelType(); t != "session" {
		newChannel.Reject(ssh.UnknownChannelType, fmt.Sprintf("unknown channel type: %s", t))
		return
	}

	// At this point, we have the opportunity to reject the client's
	// request for another logical connection
	connection, requests, err := newChannel.Accept()
	if err != nil {
		log.Printf("Could not accept channel (%s)", err)
		return
	}

	// TODO: extract into own file
	// Terminal creation code inspired by this:
	// https://github.com/antha-lang/antha/blob/master/bvendor/golang.org/x/net/http2/h2i/h2i.go
	t := terminal.NewTerminal(connection, "λ > ")
	go func() {
		for {
			line, err := t.ReadLine()
			if err == io.EOF {
				return
			}
			if err != nil {
				log.Printf("terminal.ReadLine: %v", err)
			}
			f := strings.Fields(line)
			if len(f) == 0 {
				continue
			}

			if f[0] == "exit" {
				// TODO: close session
				connection.Close()
				return
			}
			cmd := exec.Command("bash", append([]string{"-c"}, strings.Join(f[:], " "))...)
			output, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Fprintf(t, string(output))
				fmt.Fprintf(t, "%v\n", err)
				continue
			}

			fmt.Fprintf(t, string(output))
		}
	}()

	// TODO: ASCII art, version number
	t.Write([]byte("Welcome to Lambda Shell!\n"))

	go s.processRequests(t, requests)
}

// Sessions have out-of-band requests such as "shell", "pty-req" and "env"
// Good reference: https://github.com/ilowe/cmd/blob/72efdd2f2e6192e86adf67703a6f54b8bf3afc0c/sshpit/main.go
func (s *SecureServer) processRequests(t *terminal.Terminal, requests <-chan *ssh.Request) {
	var hasShell bool
	for req := range requests {
		var width, height int
		var ok bool
		switch req.Type {
		case "shell":
			if !hasShell {
				ok = true
				hasShell = true
			}
		case "exec":
			ok = true
		case "pty-req":
			width, height, ok = parsePtyReq(req.Payload)
			if ok {
				err := t.SetSize(width, height)
				ok = err == nil
			}
		case "window-change":
			width, height, ok = parseWindowChangeReq(req.Payload)
			if ok {
				err := t.SetSize(width, height)
				ok = err == nil
			}
		}

		if req.WantReply {
			req.Reply(ok, nil)
		}
	}
}
