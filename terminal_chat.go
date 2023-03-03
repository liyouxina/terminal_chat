package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"sync"
	"time"
)

var chats map[string][]string

var createChatLock sync.Mutex
var sendMsgLock sync.Mutex

func main()  {
	chats = map[string][]string{}
	r := gin.New()
	r.GET("/", index)
	r.GET("/index.html", index)
	r.GET("/chat.html", chatHtml)
	r.GET("/chat", chat)
	r.GET("/send", send)
	r.Run(":7070")
}

func chatHtml(c *gin.Context) {
	content, _ := os.ReadFile("./chat.html")
	room := c.Query("room")
	user := c.Query("user")
	kv := map[string]*string{}
	kv["room"] = &room
	kv["user"] = &user
	c.Writer.Write(templateFill(content, kv))
}

func index(c *gin.Context) {
	content, _ := os.ReadFile("./index.html")
	c.Writer.Write(content)
}

func chat(c *gin.Context) {
	roomNumber := c.Query("room")
	if isEmpty(roomNumber) {
		c.Writer.WriteString("没有这个房间号")
		return
	}

	createChatIfNotExists(roomNumber)

	contents, _ := chats[roomNumber]
	for _, content := range contents {
		c.Writer.WriteString(content)
		c.Writer.WriteString("\n")
	}
}

func send(c *gin.Context) {
	roomNumber := c.Query("room")
	if isEmpty(roomNumber) {
		c.Writer.WriteString("没有这个房间号")
		return
	}
	userNumber := c.Query("user")
	if isEmpty(userNumber) {
		c.Writer.WriteString("你是谁")
		return
	}
	content := c.Query("content")
	if isEmpty(content) {
		c.Writer.WriteString("发送消息不能为空")
		return
	}

	createChatIfNotExists(roomNumber)

	addContent(roomNumber, fmt.Sprintf("%s %s: %s</br>", time.Now().Format("01-02-2006 15:04:05"), userNumber, content))

	c.Writer.WriteString("发送成功")
}

func createChatIfNotExists(roomNumber string) {
	if isEmpty(roomNumber) {
		return
	}

	if _, exists := chats[roomNumber]; !exists {
		createChatLock.Lock()
		if _, stillExists := chats[roomNumber]; !stillExists {
			chats[roomNumber] = []string{}
		}
		createChatLock.Unlock()
	}
}

func addContent(roomNumber string, content string) {
	sendMsgLock.Lock()
	contents, _ := chats[roomNumber]
	contents = append(contents, content)
	chats[roomNumber] = contents
	sendMsgLock.Unlock()
}

func isEmpty(s string) bool {
	if len(s) == 0 {
		return true
	}
	return false
}

func templateFill(content []byte, kv map[string]*string) []byte {
	result := make([]byte, 0, 2 * len(content))
	for i := 0 ; i < len(content); i++ {
		if (i + 1 < len(content) && content[i] == '{' && content[i + 1] == '{') {
			var match []byte
			j := i + 2
			for ; j + 1 < len(content); j++ {
				if (content[j] == '}' && content[j + 1] == '}') {
					match = content[i + 2 : j]
					break
				}
			}
			if (match != nil && kv[string(match)] != nil) {
				for _, o := range []byte(*kv[string(match)]) {
					result = append(result, o)
				}
			}
			i = j + 1
		} else {
			result = append(result, content[i])
		}
	}
	return result
}