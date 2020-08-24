# go-chroot
![CI](https://github.com/kmdkuk/go-chroot/workflows/CI/badge.svg)

A template to start writing CLI easily in Go language.  
Please change "go-chroot" to your project name and use it. 

By setting your Slack webhook url in secrets.SLACK_WEBHOOK,  
you can also notify slack of test results.  

## How To Use

```shell script
$ vagrant up
$ vagrant ssh
$ cd go-chroot
# /tmp/go-chroot以下にdockerから持ってきたalpineのファイルを展開そこをrootとして，alpineの/bin/shを実行
$ sudo go run main.go
```

