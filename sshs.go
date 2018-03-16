package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/MACDfree/sshs/common"
	"github.com/MACDfree/sshs/config"
	"github.com/MACDfree/sshs/ssh"
	"github.com/fatih/color"
	"github.com/rodaine/table"

	"gopkg.in/urfave/cli.v1"
)

var (
	ip       string
	port     int
	username string
	password string
	sessions map[string]config.Session
)

func main() {
	app := cli.NewApp()

	app.Name = "sshs"
	app.Usage = "manage ssh sessions"
	app.Version = "1.0"
	app.EnableBashCompletion = true
	app.Action = func(c *cli.Context) error {
		if c.Args().First() != "" {
			var success bool
			// 读取所有配置
			sessions, success = config.ReadConfig()
			if !success {
				fmt.Println("list is empty, please execute command `sshs add` first")
				os.Exit(0)
			}
			session, ok := sessions[c.Args().First()]
			if !ok {
				fmt.Printf("do not find session named %s\n", c.Args().First())
				os.Exit(0)
			}
			ssh.SSH(session)
			return nil
		}
		fmt.Println("please execute command `sshs h` for help")
		return nil
	}

	app.Commands = []cli.Command{
		{
			Name:  "add",
			Usage: "add a ssh session to the list",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "ip, i",
					Usage: "ssh `ip` address",
				},
				cli.IntFlag{
					Name:  "port, p",
					Usage: "ssh `port`",
					Value: 22,
				},
				cli.StringFlag{
					Name:  "username, u",
					Usage: "ssh `username`",
				},
				cli.StringFlag{
					Name:  "password, w",
					Usage: "ssh `password`",
				},
			},
			Action: func(c *cli.Context) error {
				if arg := c.Args().First(); arg != "" {
					sessions, _ = config.ReadConfig()
					_, ok := sessions[arg]
					if ok {
						fmt.Printf("%s is already in the list\n", arg)
						os.Exit(0)
					}
					if c.String("ip") == "" {
						fmt.Println("--ip/-i flag not set")
						os.Exit(0)
					}
					if c.String("port") == "" {
						fmt.Println("--port/-p flag not set")
						os.Exit(0)
					}
					if c.String("username") == "" {
						fmt.Println("--username/-u flag not set")
						os.Exit(0)
					}
					sessions[arg] = config.Session{
						IP:       c.String("ip"),
						Port:     c.Int("port"),
						UserName: c.String("username"),
						Password: c.String("password"),
					}
					config.WriteConfig(sessions)
					fmt.Println("success")
					return nil
				}
				fmt.Println("alias argument not set")
				return nil
			},
		},
		{
			Name:  "rm",
			Usage: "remove a ssh session on the list",
			Action: func(c *cli.Context) error {
				if arg := c.Args().First(); arg != "" {
					var success bool
					sessions, success = config.ReadConfig()
					if !success {
						fmt.Println("list is empty, please execute command `sshs add` first")
						os.Exit(0)
					}
					_, ok := sessions[arg]
					if !ok {
						fmt.Printf("%s not found", arg)
						os.Exit(0)
					}
					delete(sessions, arg)
					config.WriteConfig(sessions)
					fmt.Println("success")
					return nil
				}
				fmt.Println("alias argument not set")
				return nil
			},
		},
		{
			Name:  "ls",
			Usage: "show session list",
			Action: func(c *cli.Context) error {
				var success bool
				sessions, success = config.ReadConfig()
				if !success {
					fmt.Println("list is empty, please execute command `sshs add` first")
					os.Exit(0)
				}
				arg := c.Args().First()
				headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
				columnFmt := color.New(color.FgYellow).SprintfFunc()
				tbl := table.New("alias", "ip", "port", "username")
				tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
				for k, v := range sessions {
					if arg != "" && !strings.Contains(k, arg) {
						continue
					}
					tbl.AddRow(k, v.IP, v.Port, v.UserName)
				}
				tbl.Print()
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	common.CheckError(err)
}

func init() {
	// 初始化设置日志格式
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
}
