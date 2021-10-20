package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"ginfra/log"
	mw "ginfra/middleware"
	"ginfra/models"
	"ginfra/protocol"
	"ginfra/utils"
)

//ClaimData 缓存到jwt的登录态信息
type ClaimData struct {
	Uid          uint64
	Identifier   string
	IdentityType int
}

//generateToken 生成登录态token
func generateToken(s *models.UserAuth) (string, error) {

	// 构造SignKey: 签名和解签名需要使用一个值
	// 构造用户claims信息(负荷)
	claimData := &ClaimData{
		Uid:          s.Uid,
		Identifier:   s.Identifier,
		IdentityType: s.IdentityType,
	}

	// 根据claims生成token对象, expires = 0 表示将从配置文件中读取默认值
	token, err := mw.GenerateToken(claimData, 0)
	if err != nil {
		return "", err
	}
	return token, nil
}

func HandleClaims(c *gin.Context, claims *utils.CustomClaims) error {
	data, err := unmarshalClaimData(claims)
	if err != nil {
		return errors.New("登录态无效")
	}

	if data.Uid == 0 {
		return errors.New("登录态无效")
	}

	log.WithGinContext(c).Debug(fmt.Sprintf("claims: Uid=%d, Identifier=%s, IdentityType=%d",
		data.Uid, data.Identifier, data.IdentityType))

	c.Set("claims", data)
	c.Set(protocol.CtxUserID, strconv.FormatUint(data.Uid, 10))
	return nil
}

//getClaimData 从登录态token获取缓存信息
func getClaimData(c *gin.Context) (*ClaimData, error) {
	claims, ok := c.MustGet("claims").(*ClaimData)
	if ok {
		return claims, nil
	}

	customclaims, ok := c.MustGet("claims").(*utils.CustomClaims)
	if !ok {
		return nil, protocol.ErrCodeInvalidClaims
	}
	data, err := unmarshalClaimData(customclaims)
	if err != nil {
		log.WithGinContext(c).Error(err.Error(), zap.String("error", protocol.ErrCodeInvalidClaims.Code))
		return nil, protocol.ErrCodeInvalidClaims
	}

	if data.Uid == 0 {
		return nil, protocol.ErrCodeInvalidClaims
	}

	log.WithGinContext(c).Debug(fmt.Sprintf("claims: Uid=%d, Identifier=%s, IdentityType=%d",
		data.Uid, data.Identifier, data.IdentityType))

	return data, nil
}

func unmarshalClaimData(claims *utils.CustomClaims) (*ClaimData, error) {
	var data ClaimData
	err := json.Unmarshal(claims.Data, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}
