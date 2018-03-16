# SSH连接管理工具

[![Build Status](https://travis-ci.org/MACDfree/sshs.svg?branch=master)](https://travis-ci.org/MACDfree/sshs)

## 起源

linux中ssh命令最多只能设置别名，无法记住密码，要么就是配置密钥。
希望可以保存远程主机的用户名和密码，快捷登录。

## 使用

登录：`sshs alias`

添加：`sshs add alias -i ip -p port -u username -w password`

删除：`sshs rm alias`

查看（支持模糊匹配）：`sshs ls [alias]`

## 存储文件位置及格式

存储文件位置：~/.sshs.yml
存储格式：

``` text
# alias ip  port    username    password
123  192.168.123.123  22  root    123456
```
