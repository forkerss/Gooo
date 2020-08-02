package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var repoURL string

func initCmd() {
	flag.StringVar(&repoURL, "repo", "none", "Git 仓库地址")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `zipGitRepo version: 0.0.1
Usage: zipGitRepo [-h] [-repo]

Options:
`)
		flag.PrintDefaults()
	}
}

func main() {
	initCmd()
	flag.Parse()
	cloneAndZip()
}

func cloneAndZip() {
	if !strings.HasSuffix(repoURL, ".git") {
		fmt.Println("[!] 不是一个仓库地址")
		return
	}

	dirs := strings.Split(repoURL, "/")
	repoDirName := dirs[len(dirs)-1]
	repoDirName = strings.Replace(repoDirName, ".git", "", 1)

	// clear old `repoDirName`.zip
	if fileExists(fmt.Sprintf("./%s.zip", repoDirName)) {
		os.RemoveAll(fmt.Sprintf("./%s.zip", repoDirName))
	}

	// git clone repository
	cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("git clone %s", repoURL))
	_, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("git clone error: %s\n", err)
		return
	}
	// defer clear repository
	defer os.RemoveAll(fmt.Sprintf("./%s", repoDirName))

	// zip 压缩加密文件
	passwd := randStringRunes(12)
	zipCmd := fmt.Sprintf("zip -r  -P %s %s.zip %s", passwd, repoDirName, repoDirName)
	cmd = exec.Command("/bin/sh", "-c", zipCmd)
	_, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("zip repository error: %s\n", err)
		return
	}
	// Succ
	fmt.Printf("Succ - file: '%s.zip' passwd: %s\n", repoDirName, passwd)
}

func fileExists(path string) bool {
	/*
	*	判断所给路径文件/文件夹是否存在
	 */
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func randStringRunes(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
