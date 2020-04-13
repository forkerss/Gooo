package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func main() {
	repos := os.Args[1:]
	for _, repoaddr := range repos {
		cloneAndZip(repoaddr)
	}
}

func cloneAndZip(repoaddr string) {
	if !strings.HasSuffix(repoaddr, ".git") {
		fmt.Println("This is not a repository.")
		return
	}
	dirs := strings.Split(repoaddr, "/")
	dirName := dirs[len(dirs)-1]
	dirName = strings.Replace(dirName, ".git", "", 1)
	// clear old `dirName`.zip
	if fileExists(fmt.Sprintf("./%s.zip", dirName)) {
		os.RemoveAll(fmt.Sprintf("./%s.zip", dirName))
	}

	// git clone repository
	gitCloneCmd := fmt.Sprintf("git clone %s", repoaddr)
	cmd := exec.Command("/bin/sh", "-c", gitCloneCmd)
	_, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("git clone error: %s\n", err)
		return
	}
	// defer clear repository
	defer os.RemoveAll(fmt.Sprintf("./%s", dirName))

	// zip 压缩加密文件
	simplePasswd := randStringRunes(12)
	zipCmd := fmt.Sprintf("zip -r  -P %s %s.zip %s", simplePasswd, dirName, dirName)
	cmd = exec.Command("/bin/sh", "-c", zipCmd)
	_, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("zip repository error: %s\n", err)
		return
	}
	// Succ
	fmt.Printf("Succ - file: '%s.zip' passwd: %s\n", dirName, simplePasswd)
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
