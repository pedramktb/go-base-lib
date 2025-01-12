package secureShell

import (
	"errors"
	"fmt"
	"net"
	"os"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// SecureShell is a struct that contains an SSH client and an SFTP client for a remote server and the username of the authenticating user.
type SecureShell struct {
	sshClient  *ssh.Client
	sftpClient *sftp.Client
	username   string
}

// New creates a new SecureShell with the given IP, username, signer and or password, and optionally additional SSH keys
// If only a signer or a password is provided, the provided method is used to authenticate.
// If both a signer and a password are provided, the password is used to authenticate and the signer's public key along with any additional SSH keys are added to the server.
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

// addKeyToServer adds the signer's public key and any additional SSH keys to the server at the given IP using the provided username and password.
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

// IP returns the IP of the remote server.
func (s *SecureShell) IP() (string, error) {
	ip, _, err := net.SplitHostPort(s.sshClient.Conn.RemoteAddr().String())
	return ip, err
}

// Run runs the given command on the remote server. Run adds "sudo" to the command if the authenticating user is not root.
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

// Upload uploads the file at the local path to the remote path on the server.
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

// Close closes the SSH and SFTP clients.
func (s *SecureShell) Close() error {
	return errors.Join(
		s.sshClient.Close(),
		s.sftpClient.Close(),
	)
}

// Username returns the username of the authenticating user.
func (s *SecureShell) Username() string {
	return s.username
}
