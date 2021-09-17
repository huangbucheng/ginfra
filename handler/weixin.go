package handler

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"

	"ginfra/config"
	"ginfra/log"
	"ginfra/protocol"
	"ginfra/utils"
)

var cfg *config.Config

func init() {
	var err error
	cfg, err = config.Parse("")
	if err != nil {
		panic(err)
	}
}

// WXCheckSignature 微信接入校验
func WXCheckSignature(c *gin.Context) {
	signature := c.Query("signature")
	timestamp := c.Query("timestamp")
	nonce := c.Query("nonce")
	echostr := c.Query("echostr")

	Token := cfg.GetString("wx.SignatureToken")
	ok := utils.CheckWxOffiAcctSignature(signature, timestamp, nonce, Token)
	if !ok {
		log.WithGinContext(c).Error("微信公众号接入校验失败!")
		return
	}

	log.WithGinContext(c).Info("微信公众号接入校验成功!")
	_, _ = c.Writer.WriteString(echostr)
}

// WXTextMsg 微信文本消息结构体
type WXTextMsg struct {
	ToUserName   string
	FromUserName string
	CreateTime   int64
	MsgType      string
	Content      string
	MsgId        int64
}

// WXMsgReceive 微信消息接收
func WXMsgReceive(c *gin.Context) {
	var textMsg WXTextMsg
	err := c.ShouldBindXML(&textMsg)
	if err != nil {
		log.WithGinContext(c).Error(fmt.Sprintf("[消息接收] - XML数据包解析失败: %v", err))
		return
	}

	log.WithGinContext(c).Info(fmt.Sprintf(
		"[消息接收] - 收到 %s 消息, 消息类型为: %s, 消息内容为: %s\n",
		textMsg.FromUserName, textMsg.MsgType, textMsg.Content))

	lb := &LiveBullet{
		Msg:       textMsg.Content,
		IsShow:    false,
		From:      textMsg.FromUserName,
		Timestamp: textMsg.CreateTime,
	}
	b, _ := json.Marshal(lb)

	utils.InsertDocuments("test-envid", "test-collection", [][]byte{b})
}

//LiveBullet 示例
type LiveBullet struct {
	Msg       string
	IsShow    bool
	From      string
	Timestamp int64
}

//JsonMsgReceiveResponse 示例
type JsonMsgReceiveResponse struct {
	DocIds []string `json:"docIds"`
}

// JsonMsgReceive 微信消息接收
func JsonMsgReceive(c *gin.Context) {
	var textMsg WXTextMsg
	err := c.ShouldBindJSON(&textMsg)
	if err != nil {
		log.WithGinContext(c).Error(fmt.Sprintf("[消息接收] - JSON数据包解析失败: %v", err))
		return
	}

	log.WithGinContext(c).Info(fmt.Sprintf(
		"[消息接收] - 收到 %s 消息, 消息类型为: %s, 消息内容为: %s\n",
		textMsg.FromUserName, textMsg.MsgType, textMsg.Content))

	lb := &LiveBullet{
		Msg:       textMsg.Content,
		IsShow:    false,
		From:      textMsg.FromUserName,
		Timestamp: textMsg.CreateTime,
	}
	b, _ := json.Marshal(lb)

	ids, err := utils.InsertDocuments("test-envid", "test-collection", [][]byte{b})
	if err != nil {
		protocol.SetErrResponse(c, err)
		return
	}

	var resp JsonMsgReceiveResponse
	resp.DocIds = ids
	protocol.SetResponse(c, &resp)
}
