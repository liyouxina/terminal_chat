package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
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
	r.GET("/chat", chat)
	r.GET("/send", send)
	r.Run(":7070")
}

func index(c *gin.Context) {

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

	addContent(roomNumber, fmt.Sprintf("%s %s: %s", time.Now().Format("01-02-2006 15:04:05"), userNumber, content))

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