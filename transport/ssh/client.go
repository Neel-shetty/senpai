package ssh

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/ssh"
)

type SSHClient struct {
	client *ssh.Client
}

type GitSession struct {
	session *ssh.Session
	stdin   io.WriteCloser
	stdout  io.Reader
	stderr  io.Reader
}

func NewSSHClient(user, host string) (*SSHClient, error) {
	key, err := os.ReadFile(os.ExpandEnv("$HOME/.ssh/id_ed25519"))
	if err != nil {
		return nil, fmt.Errorf("read private key: %w", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("parse private key: %w", err)
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", host+":22", config)
	if err != nil {
		return nil, fmt.Errorf("failed to dial ssh: %w", err)
	}

	return &SSHClient{client: conn}, nil
}

func (c *SSHClient) RunCommand(cmd string) ([]byte, error) {
	session, err := c.client.NewSession()
	if err != nil {
		return nil, fmt.Errorf("new ssh session: %w", err)
	}
	defer session.Close()

	return session.CombinedOutput(cmd)
}

func (c *SSHClient) StartGitUploadPack(repoPath string) (*GitSession, error) {
	return c.startGitSession(fmt.Sprintf("git-upload-pack '%s'", repoPath))
}

func (c *SSHClient) StartGitReceivePack(repoPath string) (*GitSession, error) {
	return c.startGitSession(fmt.Sprintf("git-receive-pack '%s'", repoPath))
}

func (c *SSHClient) startGitSession(cmd string) (*GitSession, error) {
	session, err := c.client.NewSession()
	if err != nil {
		return nil, fmt.Errorf("new ssh session: %w", err)
	}

	stdin, err := session.StdinPipe()
	if err != nil {
		session.Close()
		return nil, fmt.Errorf("get stdin pipe: %w", err)
	}

	stdout, err := session.StdoutPipe()
	if err != nil {
		session.Close()
		return nil, fmt.Errorf("get stdout pipe: %w", err)
	}

	stderr, err := session.StderrPipe()
	if err != nil {
		session.Close()
		return nil, fmt.Errorf("get stderr pipe: %w", err)
	}

	if err := session.Start(cmd); err != nil {
		session.Close()
		return nil, fmt.Errorf("start command: %w", err)
	}

	return &GitSession{
		session: session,
		stdin:   stdin,
		stdout:  stdout,
		stderr:  stderr,
	}, nil
}

func (s *GitSession) Read(p []byte) (int, error) {
	return s.stdout.Read(p)
}

func (s *GitSession) Write(p []byte) (int, error) {
	return s.stdin.Write(p)
}

func (s *GitSession) GetStdout() io.Reader {
	return s.stdout
}

func (s *GitSession) GetStderr() io.Reader {
	return s.stderr
}

func (s *GitSession) GetStdin() io.WriteCloser {
	return s.stdin
}

func (s *GitSession) CloseStdin() error {
	return s.stdin.Close()
}

func (s *GitSession) Close() error {
	s.stdin.Close()
	err := s.session.Wait()
	s.session.Close()
	return err
}

func (c *SSHClient) Close() error {
	return c.client.Close()
}
