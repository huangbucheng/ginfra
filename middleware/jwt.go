package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"

	"ginfra/config"
	"ginfra/errcode"
	"ginfra/log"
	"ginfra/protocol"
	"ginfra/utils"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	RS256PublicKey  []byte = []byte{}
	RS256PrivateKey []byte = []byte{}
	JWTExpires      int64
	JWTIssuer       string
	HeaderTokenName string
	CookieTokenName string
)

func init() {
	cfg, err := config.Parse("")
	if err != nil {
		panic(err)
	}

	// init
	RS256KeyDir := cfg.GetString("jwt.RS256KeyDir")
	privateKeyFile := filepath.Join(RS256KeyDir, "rs256.key")
	RS256PrivateKey, err = ioutil.ReadFile(privateKeyFile)
	if utils.Exists(privateKeyFile) && err != nil {
		panic(fmt.Errorf("read private key file %s error:%s", privateKeyFile, err.Error()))
	}

	publicKeyFile := filepath.Join(RS256KeyDir, "rs256.key.pub")
	RS256PublicKey, err = ioutil.ReadFile(publicKeyFile)
	if utils.Exists(publicKeyFile) && err != nil {
		panic(fmt.Errorf("read public key file %s error:%s", publicKeyFile, err.Error()))
	}

	JWTExpires = cfg.GetInt64("jwt.jwtexpires")
	JWTIssuer = cfg.GetString("jwt.jwtissuer")
	cfg.SetDefault("jwt.headername", "token")
	cfg.SetDefault("jwt.cookiename", "token")
	HeaderTokenName = cfg.GetString("jwt.headername")
	CookieTokenName = cfg.GetString("jwt.cookiename")
}

type HandleClaimFunc func(c *gin.Context, claims *utils.CustomClaims) error

// JWTAuth 中间件，检查token
func JWTAuth(claimHandler HandleClaimFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		var token string
		token = c.Request.Header.Get(HeaderTokenName)
		if token == "" {
			token, err = c.Cookie(CookieTokenName)
			if err != nil {
				log.WithGinContext(c).Error("JWTAuth no token")
				protocol.SetErrResponse(c, errcode.NewCustomError(errcode.ErrNoAuthToken, "no auth token"))
				c.Abort()
				return
			}
		}

		// 解析token中包含的相关信息
		claims, err := ParseToken(token)
		if err != nil {
			// token过期
			log.WithGinContext(c).Error("JWTAuth ParseJWTTokenWithRS256 fail", zap.String("error", err.Error()))
			if _, ok := err.(*jwt.TokenExpiredError); ok {
				protocol.SetErrResponse(c, errcode.NewCustomError(errcode.ErrExpiredAuthToken, "expired auth token"))
				c.Abort()
				return
			}
			protocol.SetErrResponse(c, errcode.NewCustomError(errcode.ErrInvalidAuthToken, "invalid auth token"))
			c.Abort()
			return
		}

		// 解析到具体的claims相关信息
		//c.Set("claims", claims)
		err = claimHandler(c, claims)
		if err != nil {
			log.WithGinContext(c).Error("JWTAuth HandleClaimFunc exception", zap.String("error", err.Error()))
			protocol.SetErrResponse(c, errcode.NewCustomError(errcode.ErrInvalidJWTClaims, err.Error()))
			c.Abort()
			return
		}
	}
}

//GenerateToken 生成登录态token
func GenerateToken(claims interface{}, expires int64) (string, error) {
	b, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	// 构造用户claims信息(负荷)

	// 根据claims生成token对象
	token, err := utils.CreateJWTTokenWithRS256(
		RS256PrivateKey,
		NewCustomClaims(b, expires),
	)
	if err != nil {
		return "", err
	}
	return token, nil
}

//ParseToken 解析登录态Token
func ParseToken(token string) (claims *utils.CustomClaims, err error) {
	return utils.ParseJWTTokenWithRS256(RS256PublicKey, token)
}

//GenerateSignature 生成签名串
func GenerateSignature(b []byte, expires int64) (string, error) {
	sig, err := utils.CreateJWTTokenWithRS256(
		RS256PrivateKey,
		NewCustomClaims(b, expires),
	)
	if err != nil {
		return "", err
	}
	return sig, nil
}

//VerifySignature 校验签名串
func VerifySignature(sig string) ([]byte, error) {
	claims, err := utils.ParseJWTTokenWithRS256(RS256PublicKey, sig)
	if err != nil {
		return []byte{}, err
	}
	return claims.Data, nil
}

// JWT基本数据结构
// 签名的signkey
type JWT struct {
	SigningKey []byte
}

// 初始化JWT实例
func NewJWT(key []byte) *JWT {
	return &JWT{key}
}

// 新建一个CustomClaims
func NewCustomClaims(data []byte, expires int64) *utils.CustomClaims {
	if expires == 0 {
		expires = JWTExpires
	}
	return &utils.CustomClaims{
		Data: data,
		StandardClaims: jwt.StandardClaims{
			NotBefore: jwt.At(time.Now().Add(-1 * time.Hour)),                       // 签名生效时间
			ExpiresAt: jwt.At(time.Now().Add(time.Duration(expires) * time.Second)), // 签名过期时间
			Issuer:    JWTIssuer,                                                    // 签名颁发者
		},
	}
}

// 创建Token(基于用户的基本信息claims)
// 使用HS256算法进行token生成
// 使用用户基本信息claims以及签名key(signkey)生成token
func (j *JWT) CreateToken(claims utils.CustomClaims) (string, error) {
	// https://gowalker.org/github.com/dgrijalva/jwt-go#Token
	// 返回一个token的结构体指针
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// token解析
// Couldn't handle this token:
func (j *JWT) ParserToken(tokenString string) (*utils.CustomClaims, error) {
	// https://gowalker.org/github.com/dgrijalva/jwt-go#ParseWithClaims
	// 输入用户自定义的Claims结构体对象,token,以及自定义函数来解析token字符串为jwt的Token结构体指针
	// Keyfunc是匿名函数类型: type Keyfunc func(*Token) (interface{}, error)
	// func ParseWithClaims(tokenString string, claims Claims, keyFunc Keyfunc) (*Token, error) {}
	token, err := jwt.ParseWithClaims(tokenString, &utils.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})

	if err != nil {
		// https://gowalker.org/github.com/dgrijalva/jwt-go#ValidationError
		// jwt.ValidationError 是一个无效token的错误结构
		//if _, ok := err.(*jwt.TokenExpiredError); ok {
		//	return nil, ErrTokenExpired
		//} else {
		//	return nil, ErrTokenInvalid
		//}
		return nil, err
	}

	// 将token中的claims信息解析出来和用户原始数据进行校验
	// 做以下类型断言，将token.Claims转换成具体用户自定义的Claims结构体
	if claims, ok := token.Claims.(*utils.CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("Token无效")

}

// 更新Token
func (j *JWT) UpdateToken(tokenString string) (string, error) {
	// TimeFunc为一个默认值是time.Now的当前时间变量,用来解析token后进行过期时间验证
	// 可以使用其他的时间值来覆盖
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}

	// 拿到token基础数据
	token, err := jwt.ParseWithClaims(tokenString, &utils.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil

	})

	// 校验token当前还有效
	if claims, ok := token.Claims.(*utils.CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		// 修改Claims的过期时间(int64)
		// https://gowalker.org/github.com/dgrijalva/jwt-go#StandardClaims
		claims.StandardClaims.ExpiresAt = jwt.At(time.Now().Add(1 * time.Hour))
		return j.CreateToken(*claims)
	}
	return "", fmt.Errorf("token获取失败:%v", err)
}
