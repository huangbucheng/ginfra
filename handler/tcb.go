package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"time"

	"ginfra/config"
	"ginfra/errcode"
	"ginfra/log"
	"ginfra/protocol"
	"ginfra/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//GetTicketRequest 获取TCB ticket的请求参数
type GetTicketRequest struct {
	EnvID string `binding:"required,min=1"`
}

//GetTicketResponse 获取TCB ticket的响应参数
type GetTicketResponse struct {
	Ticket string
}

type TcbPrivateKeyData struct {
	PrivateKeyID string `json:"private_key_id"`
	PrivateKey   string `json:"private_key"`
	EnvID        string `json:"env_id"`
}

//GetTicket 获取TCB ticket
func GetTicket(c *gin.Context) {
	var req GetTicketRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.WithGinContext(c).Error(err.Error(), zap.String("error", errcode.ErrInvalidParam))
		protocol.SetErrResponse(c, protocol.ErrCodeInvalidParameter)
		return
	}

	claims, err := getClaimDataFromContext(c)
	if err != nil {
		protocol.SetErrResponse(c, err)
		return
	}

	privateKeyData, err := getPrivateKey(req.EnvID)
	if err != nil {
		log.WithGinContext(c).Error(err.Error())
		protocol.SetErrResponse(c,
			errcode.NewCustomError(errcode.ErrCodeInternalError,
				"暂不支持该云开发环境Ticket签发"))
		return
	}

	data := make(map[string]interface{})
	data["uid"] = strconv.FormatUint(claims.Uid, 10)
	data["env"] = req.EnvID
	data["iat"] = time.Now().UnixNano() / 1000000                        // Token颁发时间
	data["exp"] = time.Now().Add(10*60*time.Second).UnixNano() / 1000000 // Token过期时间
	data["alg"] = "RS256"
	data["refresh"] = 3600 * 1000
	data["expire"] = time.Now().Add(7*24*time.Hour).UnixNano() / 1000000

	ticket, err := utils.CreateJWTTokenFromMapWithRS256(
		[]byte(privateKeyData.PrivateKey),
		data,
	)
	if err != nil {
		protocol.SetErrResponse(c,
			errcode.NewCustomError(errcode.ErrCodeInternalError,
				fmt.Sprintf("create ticket error:%s", err.Error())))
		return
	}

	ticket = fmt.Sprintf("%s/@@/%s", privateKeyData.PrivateKeyID, ticket)
	var resp GetTicketResponse
	resp.Ticket = ticket

	cfg, _ := config.Parse("")
	c.SetCookie("ticket", ticket, 24*3600, "/", cfg.GetString("jwt.domain"), false, true)
	protocol.SetResponse(c, &resp)
}

func getPrivateKey(envid string) (*TcbPrivateKeyData, error) {
	cfg, _ := config.Parse("")

	privateKeyDir := cfg.GetString("tcb.PrivateKeyDir")
	privateKeyFile := filepath.Join(privateKeyDir, envid)
	privateKeyContent, err := ioutil.ReadFile(privateKeyFile)
	if !utils.Exists(privateKeyFile) || err != nil {
		return nil, fmt.Errorf("private key nonexist:%s", privateKeyFile)
	}

	var privateKeyData TcbPrivateKeyData
	err = json.Unmarshal(privateKeyContent, &privateKeyData)
	if err != nil {
		return nil, fmt.Errorf("unmarshal private key err:%s", err.Error())
	}
	if privateKeyData.EnvID != envid {
		return nil, fmt.Errorf("private key's env id mismatch:%s vs %s",
			privateKeyData.EnvID, envid)
	}

	return &privateKeyData, nil
}
