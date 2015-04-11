package sshclient

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
)

type SshClient struct {
	User     string
	Server   string
	Key      string
	Port     string
	Password string
	AuthType string
	sclient  *ssh.Client
}

func NewSshClient(server, user string) *SshClient {
	sc := SshClient{
		User:     user,
		Server:   server,
		Port:     "22",
		AuthType: "PublicKey",
		Key:      "/.ssh/id_rsa",
	}
	return &sc
}

func getKeyFile(keypath string) (ssh.Signer, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}

	file := usr.HomeDir + keypath
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	pubkey, err := ssh.ParsePrivateKey(buf)
	if err != nil {
		return nil, err
	}

	return pubkey, nil
}

func (client *SshClient) connect() error {
	pubkey, err := getKeyFile(client.Key)
	if err != nil {
		return err
	}

	config := &ssh.ClientConfig{User: client.User}
	if client.AuthType == "Password" {
		config.Auth = []ssh.AuthMethod{ssh.Password(client.Password)}
	} else {
		config.Auth = []ssh.AuthMethod{ssh.PublicKeys(pubkey)}
	}

	sc, err := ssh.Dial("tcp", client.Server+":"+client.Port, config)
	if err != nil {
		return err
	}
	client.sclient = sc
	return nil
}

func (client *SshClient) session() (*ssh.Session, error) {
	if client.sclient == nil {
		if err := client.connect(); err != nil {
			return nil, err
		}
	}

	sc := client.sclient
	session, err := sc.NewSession()
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (client *SshClient) Run(command string) (string, error) {
	session, err := client.session()
	if err != nil {
		return "", err
	}
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b
	err = session.Run(command)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}

func (client *SshClient) Scp(sourceFile string) error {
	session, err := client.session()
	if err != nil {
		return err
	}
	defer session.Close()

	targetFile := filepath.Base(sourceFile)
	src, srcErr := os.Open(sourceFile)

	if srcErr != nil {
		return srcErr
	}

	srcStat, statErr := src.Stat()
	if statErr != nil {
		return statErr
	}

	go func() {
		w, _ := session.StdinPipe()
		fmt.Fprintln(w, "C0644", srcStat.Size(), targetFile)
		if srcStat.Size() > 0 {
			io.Copy(w, src)
			fmt.Fprint(w, "\x00")
			w.Close()
		} else {
			fmt.Fprint(w, "\x00")
			w.Close()
		}
	}()

	if err := session.Run(fmt.Sprintf("scp -t %s", targetFile)); err != nil {
		return err
	}

	return nil
}

func (client *SshClient) Close() {
	if client.sclient != nil {
		client.sclient.Close()
	}
}
