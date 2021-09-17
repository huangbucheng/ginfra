package handler

import (
	"encoding/json"

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

//GenerateToken 生成登录态token
func GenerateToken(s *models.UserAuth) (string, error) {

	// 构造SignKey: 签名和解签名需要使用一个值
	// 构造用户claims信息(负荷)
	claimData := &ClaimData{
		Uid:          s.Uid,
		Identifier:   s.Identifier,
		IdentityType: s.IdentityType,
	}

	// 根据claims生成token对象
	token, err := mw.GenerateToken(claimData, 0)
	if err != nil {
		return "", err
	}
	return token, nil
}

//GetClaimData 从登录态token获取缓存信息
func GetClaimData(claims *utils.CustomClaims) (*ClaimData, error) {
	var data ClaimData
	err := json.Unmarshal(claims.Data, &data)
	if err != nil {
		return nil, protocol.ErrCodeInvalidClaims
	}
	if data.Uid == 0 {
		return nil, protocol.ErrCodeInvalidClaims
	}
	return &data, nil
}

//getClaimDataFromContext 从登录态token获取缓存信息
func getClaimDataFromContext(c *gin.Context) (*ClaimData, error) {
	claims, ok := c.MustGet("claims").(*utils.CustomClaims)
	if !ok {
		return nil, protocol.ErrCodeInvalidClaims
	}

	data, err := GetClaimData(claims)
	if err != nil {
		log.WithGinContext(c).Error(err.Error(), zap.String("error", protocol.ErrCodeInvalidClaims.Code))
	}

	return data, err
}
