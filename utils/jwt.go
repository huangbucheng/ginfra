package utils

import (
	"errors"

	"github.com/dgrijalva/jwt-go/v4"
)

// 定义载荷
type CustomClaims struct {
	Data []byte `json:"data"`
	// StandardClaims结构体实现了Claims接口(Valid()函数)
	jwt.StandardClaims
}

// CreateJWTTokenWithRS256 生成一个RS256验证的Token
// Token里面包括的值，可以自己根据情况添加，
func CreateJWTTokenWithRS256(privateKey []byte, claims *CustomClaims) (
	tokenStr string, err error) {
	_privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(_privateKey)
}

// ParseJWTTokenWithRS256 获取Token的主题（也可以更改获取其他值）
// 参数tokenStr指的是 从客户端传来的待验证Token
// 验证Token过程中，如果Token生成过程中，指定了iat与exp参数值，将会自动根据时间戳进行时间验证
func ParseJWTTokenWithRS256(publicKey []byte, tokenStr string) (
	claims *CustomClaims, err error) {
	_publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKey)
	if err != nil {
		return nil, err
	}

	// 基于公钥验证Token合法性
	//token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 基于JWT的第一部分中的alg字段值进行一次验证
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("验证Token的加密类型错误")
		}
		return _publicKey, nil
	}, jwt.WithoutAudienceValidation())
	if err != nil {
		//	if token != nil {
		//		b, _ := json.Marshal(token)
		//		fmt.Println("--->", string(b))
		//	}
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("Token无效")
}

// CreateJWTTokenFromMapWithRS256 生成一个RS256验证的Token
// Token里面包括的值，可以自己根据情况添加，
func CreateJWTTokenFromMapWithRS256(privateKey []byte, claims map[string]interface{}) (
	tokenStr string, err error) {
	_privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		return "", err
	}

	fullclaims := jwt.MapClaims{}
	for key, value := range claims {
		fullclaims[key] = value
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, fullclaims)
	return token.SignedString(_privateKey)
}

// ParseJWTTokenFromMapWithRS256 获取Token的主题（也可以更改获取其他值）
// 参数tokenStr指的是 从客户端传来的待验证Token
// 验证Token过程中，如果Token生成过程中，指定了iat与exp参数值，将会自动根据时间戳进行时间验证
func ParseJWTTokenFromMapWithRS256(publicKey []byte, tokenStr string) (
	result map[string]interface{}, err error) {
	_publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKey)
	if err != nil {
		return nil, err
	}

	// 基于公钥验证Token合法性
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// 基于JWT的第一部分中的alg字段值进行一次验证
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("验证Token的加密类型错误")
		}
		return _publicKey, nil
	}, jwt.WithoutAudienceValidation())
	if err != nil {
		//	if token != nil {
		//		b, _ := json.Marshal(token)
		//		fmt.Println("--->", string(b))
		//	}
		return nil, err
	}

	result = make(map[string]interface{})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		for key, value := range claims {
			result[key] = value
		}
		return result, nil
	}

	return nil, errors.New("Token无效")
}
