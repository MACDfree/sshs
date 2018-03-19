package ssh

import (
	"io/ioutil"
	"net"
	"os"
	"strconv"

	"github.com/MACDfree/sshs/common"
	"github.com/MACDfree/sshs/config"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

// SSH 建立连接
func SSH(session config.Session) {
	var clientConfig *ssh.ClientConfig
	if session.Password != "" {
		clientConfig = &ssh.ClientConfig{
			User: session.UserName,
			Auth: []ssh.AuthMethod{
				ssh.Password(session.Password),
			},
			HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				return nil
			},
		}
	} else {
		key, err := ioutil.ReadFile(common.HomePath() + "/.ssh/id_rsa")
		common.CheckError(err)
		signer, err := ssh.ParsePrivateKey(key)
		common.CheckError(err)
		clientConfig = &ssh.ClientConfig{
			User: session.UserName,
			Auth: []ssh.AuthMethod{
				ssh.PublicKeys(signer),
			},
			HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				return nil
			},
		}
	}

	// int转string需要使用strconv.Itoa，不能直接使用string()强转
	client, err := ssh.Dial("tcp", session.IP+":"+strconv.Itoa(session.Port), clientConfig)
	common.CheckError(err)
	defer client.Close()
	sshSession, err := client.NewSession()
	common.CheckError(err)
	defer sshSession.Close()

	fd := int(os.Stdin.Fd())
	oldState, err := terminal.MakeRaw(fd)
	common.CheckError(err)
	sshSession.Stdout = os.Stdout
	sshSession.Stderr = os.Stderr
	sshSession.Stdin = os.Stdin

	termWidth, termHeight, err := terminal.GetSize(fd)
	common.CheckError(err)
	defer terminal.Restore(fd, oldState)

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	if err = sshSession.RequestPty("xterm-256color", termHeight, termWidth, modes); err != nil {
		common.CheckError(err)
	}
	err = sshSession.Shell()
	common.CheckError(err)
	err = sshSession.Wait()
	common.CheckError(err)
}
