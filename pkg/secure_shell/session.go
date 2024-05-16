package secureShell

type session interface {
	CombinedOutput(cmd string) ([]byte, error)
	Close() error
}

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

type sudoSession struct {
	session session
}

func newSudoSession(s session) session {
	return &sudoSession{session: s}
}

func (s *sudoSession) CombinedOutput(cmd string) ([]byte, error) {
	return s.session.CombinedOutput("sudo " + cmd)
}

func (s *sudoSession) Close() error {
	return s.session.Close()
}
