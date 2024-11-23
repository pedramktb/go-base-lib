package secureShell

import (
	"errors"
	"fmt"
	"net"
	"os"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type SecureShell struct {
	sshClient  *ssh.Client
	sftpClient *sftp.Client
	username   string
}

func New(ip, username string, signer ssh.Signer, password string, additionalSSHKeys ...string) (*SecureShell, error) {
	var config *ssh.ClientConfig
	if signer != nil {
		if password != "" {
			err := addKeyToServer(ip, username, signer, password, additionalSSHKeys...)
			if err != nil {
				return nil, err
			}
		}
		config = &ssh.ClientConfig{
			User:            username,
			Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}
	} else {
		config = &ssh.ClientConfig{
			User:            username,
			Auth:            []ssh.AuthMethod{ssh.Password(password)},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}
	}

	sshClient, err := ssh.Dial("tcp", ip+":22", config)
	if err != nil {
		return nil, err
	}

	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		return nil, err
	}

	return &SecureShell{
		sshClient:  sshClient,
		sftpClient: sftpClient,
		username:   username,
	}, nil
}

func addKeyToServer(ip, username string, signer ssh.Signer, tempPassword string, additionalSSHKeys ...string) error {
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(tempPassword),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	sshClient, err := ssh.Dial("tcp", ip+":22", config)
	if err != nil {
		return err
	}
	defer sshClient.Close()

	sess, err := sshClient.NewSession()
	if err != nil {
		return err
	}
	defer sess.Close()

	keys := string(ssh.MarshalAuthorizedKey(signer.PublicKey()))
	for i := range additionalSSHKeys {
		keys += additionalSSHKeys[i] + "\n"
	}

	out, err := sess.CombinedOutput(
		`mkdir -p ~/.ssh;
		echo "` + keys + `" >> ~/.ssh/authorized_keys;`,
	)
	if err != nil {
		return fmt.Errorf("%w: %s", err, string(out))
	}

	return nil
}

func (s *SecureShell) IP() (string, error) {
	ip, _, err := net.SplitHostPort(s.sshClient.Conn.RemoteAddr().String())
	return ip, err
}

func (s *SecureShell) Run(cmd string) error {
	sess, err := s.newSession()
	if err != nil {
		return err
	}
	defer sess.Close()

	out, err := sess.CombinedOutput(cmd)
	if err != nil {
		return fmt.Errorf("%w: %s", err, string(out))
	}

	return nil
}

func (s *SecureShell) Upload(localPath, remotePath string) error {
	src, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := s.sftpClient.Create(remotePath)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = dst.ReadFrom(src)
	if err != nil {
		return err
	}

	return nil
}

func (s *SecureShell) Close() error {
	return errors.Join(
		s.sshClient.Close(),
		s.sftpClient.Close(),
	)
}

func (s *SecureShell) Username() string {
	return s.username
}
