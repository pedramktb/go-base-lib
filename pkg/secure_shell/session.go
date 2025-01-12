package secureShell

// session is an interface that ssh.Session and sudoSession implement
type session interface {
	CombinedOutput(cmd string) ([]byte, error)
	Close() error
}

// newSession creates a new session with the remote server.
// If the authenticating user is not root, newSession returns a sudoSession that wraps the ssh.Session.
func (s *SecureShell) newSession() (session, error) {
	sess, err := s.sshClient.NewSession()
	if err != nil {
		return nil, err
	}

	if s.username == "root" {
		return sess, nil
	} else {
		return newSudoSession(sess), nil
	}
}

// sudoSession is a struct that wraps an ssh.Session and runs commands with "sudo"
type sudoSession struct {
	session session
}

// newSudoSession creates a new sudoSession that wraps the given session
func newSudoSession(s session) session {
	return &sudoSession{session: s}
}

// CombinedOutput runs the given command with "sudo" on the remote server and returns the combined output
func (s *sudoSession) CombinedOutput(cmd string) ([]byte, error) {
	return s.session.CombinedOutput("sudo " + cmd)
}

// Close closes the sudoSession
func (s *sudoSession) Close() error {
	return s.session.Close()
}
