package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/user"
	"regexp"
	"strings"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	home     string
	path     string
	ip       string
	port     int
	username string
	password string
)

func main() {
	if len(os.Args) == 2 {
		switch os.Args[1] {
		case "ls":
			lsSession("")
			os.Exit(0)
		default:
			openSession(os.Args[1])
		}
	} else if len(os.Args) == 3 {
		switch os.Args[1] {
		case "rm":
			rmSession(os.Args[2])
			os.Exit(0)
		case "ls":
			lsSession(os.Args[2])
			os.Exit(0)
		default:
			fmt.Println("命令格式错误！")
			os.Exit(0)
		}
	} else {
		flag.CommandLine.Parse(os.Args[3:])
		if ip == "" {
			fmt.Println("远程主机IP不能为空")
			os.Exit(0)
		}
		if port == 0 {
			fmt.Println("远程主机端口号不能为空")
			os.Exit(0)
		}
		if username == "" {
			fmt.Println("远程主机用户名不能为空")
			os.Exit(0)
		}
		addSession(os.Args[2], ip, port, username, password)
	}
}

func init() {
	// 初始化设置日志格式
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)

	u, err := user.Current()
	if err != nil {
		log.Fatalf("获取当前用户出错，错误信息：%s\n", err)
	}
	home = u.HomeDir
	// 设置配置文件路径
	path = u.HomeDir + "/.sshs"

	flag.StringVar(&ip, "i", "", "远程主机IP")
	flag.IntVar(&port, "p", 0, "远程主机端口号")
	flag.StringVar(&username, "u", "", "远程主机用户名")
	flag.StringVar(&password, "w", "", "远程主机密码")
}

func openSession(alias string) {
	if !checkFileIsExist(path) {
		fmt.Printf("配置文件不存在，[%s]", path)
		os.Exit(0)
	}
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("读取文件内容出错，错误信息：%s\n", err)
	}
	r := regexp.MustCompile("(?m)^" + alias + "(\\s+\\S+){3,4}\\s*$")
	text := r.FindAllString(string(content), -1)
	if len(text) > 1 {
		fmt.Printf("配置文件中存在相同的别名[%s]", alias)
		os.Exit(0)
	}
	if len(text) < 1 {
		fmt.Printf("配置文件中未找到此别名[%s]", alias)
		os.Exit(0)
	}

	line := text[0]
	line = strings.Trim(line, " ")
	//log.Println(line)
	parms := regexp.MustCompile("\\s+").Split(line, -1)
	//log.Println(parms)
	var password string
	if len(parms) == 5 {
		password = parms[4]
	} else {
		password = ""
	}
	newSSH(parms[1]+":"+parms[2], parms[3], password)
}

func rmSession(alias string) {
	if !checkFileIsExist(path) {
		fmt.Printf("配置文件不存在，[%s]", path)
		os.Exit(0)
	}
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("读取文件内容出错，错误信息：%s\n", err)
	}
	r := regexp.MustCompile("(?m)^" + alias + "(\\s+\\S+){3,4}$")
	lines := r.FindAllString(string(content), -1)
	if len(lines) > 1 {
		fmt.Println("匹配到多条数据，请确认")
		for _, str := range lines {
			fmt.Println(str)
		}
		os.Exit(0)
	}

	if len(lines) < 1 {
		fmt.Println("无匹配项")
		os.Exit(0)
	}

	text := strings.Replace(string(content), lines[0]+"\n", "", -1)
	err = ioutil.WriteFile(path, []byte(text), 0600)
	if err != nil {
		log.Fatalf("写入文件出错，错误信息：%s\n", err)
	}
}

func lsSession(alias string) {
	if !checkFileIsExist(path) {
		fmt.Printf("配置文件不存在，[%s]", path)
		os.Exit(0)
	}
	content, err := ioutil.ReadFile(path)
	//log.Print(content)
	if err != nil {
		log.Fatalf("读取文件内容出错，错误信息：%s\n", err)
	}
	var regstr string
	if alias == "" {
		regstr = "(?m)^[^#\n]\\S*?(\\s+\\S+){3}"
	} else {
		regstr = "(?m)^" + alias + "\\S*?(\\s+\\S+){3}"
	}
	r := regexp.MustCompile(regstr)
	strs := r.FindAllString(string(content), -1)
	//log.Println(r)
	fmt.Printf("%s\t%s\t%s\t%s\n", "alias", "ip", "port", "username")
	for _, str := range strs {
		fmt.Println(str)
	}
	os.Exit(0)
}

func addSession(alias, ip string, port int, username, password string) {
	if checkExist(alias) {
		fmt.Printf("别名 [%s] 已存在\n", alias)
		os.Exit(0)
	}
	var str string
	if password == "" {
		str = fmt.Sprintf("%s\t%s\t%d\t%s\n", alias, ip, port, username)
	} else {
		str = fmt.Sprintf("%s\t%s\t%d\t%s\t%s\n", alias, ip, port, username, password)
	}
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		log.Fatalf("打开配置文件出错，错误信息：%s\n", err)
	}
	defer file.Close()
	_, err = io.WriteString(file, str)
	if err != nil {
		log.Fatalf("写入配置文件出错，错误信息：%s\n", err)
	}
}

// 判断是否已存在此别名
func checkExist(alias string) bool {
	if !checkFileIsExist(path) {
		return false
	}
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("读取文件内容出错，错误信息：%s\n", err)
	}
	matched, err := regexp.MatchString("(?m)^"+alias+"(\\s+\\S+){3,4}$", string(content))
	if err != nil {
		log.Fatalf("判断是否存在出错，错误信息：%s\n", err)
	}
	return matched
}

func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func newSSH(ipport, username, password string) {
	var config *ssh.ClientConfig
	if password != "" {
		config = &ssh.ClientConfig{
			User: username,
			Auth: []ssh.AuthMethod{
				ssh.Password(password),
			},
			HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				return nil
			},
		}
	} else {
		key, err := ioutil.ReadFile(home + "/.ssh/id_rsa")
		if err != nil {
			log.Fatalf("读取私钥出错，错误信息：%s\n", err)
		}
		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			log.Fatalf("转换私钥出错，错误信息：%s\n", err)
		}
		config = &ssh.ClientConfig{
			User: username,
			Auth: []ssh.AuthMethod{
				ssh.PublicKeys(signer),
			},
			HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				return nil
			},
		}
	}

	client, err := ssh.Dial("tcp", ipport, config)
	if err != nil {
		log.Fatalf("获取ssh连接出错，错误信息：%s\n", err)
		return
	}
	defer client.Close()
	session, err := client.NewSession()
	if err != nil {
		log.Fatalf("获取session出错，错误信息：%s\n", err)
		return
	}
	defer session.Close()

	fd := int(os.Stdin.Fd())
	oldState, err := terminal.MakeRaw(fd)
	if err != nil {
		log.Fatalf("获取ssh连接出错，错误信息：%s\n", err)
		return
	}
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	termWidth, termHeight, err := terminal.GetSize(fd)
	if err != nil {
		log.Fatalf("获取当前终端大小出错，错误信息：%s\n", err)
		return
	}
	defer terminal.Restore(fd, oldState)

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	if err = session.RequestPty("xterm-256color", termHeight, termWidth, modes); err != nil {
		log.Fatalf("请求远程主机连接pty出错，错误信息：%s\n", err)
		return
	}
	err = session.Shell()
	if err != nil {
		log.Fatalf("创建shell出错，错误信息：%s\n", err)
		return
	}
	err = session.Wait()
	if err != nil {
		log.Fatalf("等待远程命令退出出错，错误信息：%s\n", err)
		return
	}
}
