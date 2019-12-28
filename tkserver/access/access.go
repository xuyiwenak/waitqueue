package access

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var (
	// tokenExpiredDate app token过期日期 30天
	tokenExpiredDate = 3600 * 24 * 30 * time.Second
	//
	generateTokenSecret = "abcdhudads"

)
func MakeAccessToken(seed int64) (ret string, err error) {


	m, err := createTokenClaims(seed)
	if err != nil {
		return "", fmt.Errorf("[MakeAccessToken] 创建token Claim 失败，err: %s", err)
	}

	// 创建
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, m)
	if tokenStr, err := token.SignedString([]byte(generateTokenSecret));err!=nil{
		return "", fmt.Errorf("[MakeAccessToken] 创建token失败，err: %s", err)
	}else{
		return tokenStr, nil
	}
}

func createTokenClaims(timeSeed int64) (m *jwt.StandardClaims, err error) {

	now := time.Now()
	m = &jwt.StandardClaims{
		ExpiresAt: now.Add(tokenExpiredDate).Unix(),
		NotBefore: now.Unix(),
		Id:        string(timeSeed),
		IssuedAt:  now.Unix(),
		Issuer:    "micro.bambooRat",
		Subject:   string(timeSeed),
	}
	return
}

