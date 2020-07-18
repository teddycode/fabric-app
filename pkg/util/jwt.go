package util

import (
	"github.com/fabric-app/models"
	"github.com/fabric-app/pkg/setting"
	"github.com/dgrijalva/jwt-go"
	"strconv"
	"strings"
	"time"
)

type Claims struct {
	jwt.StandardClaims
}

//生成令牌
func GenerateToken(user models.User) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(3 * time.Hour)
	claims := Claims{
		jwt.StandardClaims{
			Audience:  user.UserName,         // 受众
			ExpiresAt: expireTime.Unix(),     // 失效时间
			Id:        strconv.Itoa(user.ID), // 编号
			IssuedAt:  time.Now().Unix(),     // 签发时间
			Issuer:    "github.com/fabric-app",           // 签发人
			NotBefore: time.Now().Unix(),     // 生效时间
			Subject:   "login",               // 场景
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtSecret := []byte(setting.JwtSecret + user.Secret)
	return tokenClaims.SignedString(jwtSecret)
}

//解析令牌
func ParseToken(token string) (*Claims, error) {

	u, e := tokenUser(token)
	if e != nil {
		return nil, e
	}
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(setting.JwtSecret + u.Secret), nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}

//刷新令牌
func RefreshToken(tokenString string) (string, error) {
	u, e := tokenUser(tokenString)

	if e != nil {
		return "", e
	}

	tokenClaims, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(setting.JwtSecret + u.Secret), nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			id, err := strconv.Atoi(claims.Id)
			if err != nil {
				return "", err
			}
			user, userError := models.FindUserById(id)
			if userError != nil {
				return "", err
			}
			nowTime := time.Now()
			expireTime := nowTime.Add(1 * time.Hour)
			claims := Claims{
				jwt.StandardClaims{
					Audience:  user.Email,            // 受众
					ExpiresAt: expireTime.Unix(),     // 失效时间
					Id:        strconv.Itoa(user.ID), // 编号
					IssuedAt:  time.Now().Unix(),     // 签发时间
					Issuer:    "github.com/fabric-app",           // 签发人
					NotBefore: time.Now().Unix(),     // 生效时间
					Subject:   "refresh",             // 场景
				},
			}
			tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			return tokenClaims.SignedString([]byte(setting.JwtSecret + u.Secret))
		}
	}
	return "", err
}

//根据token 获取user
func tokenUser(token string) (models.User, error) {
	user := models.User{}
	payload := strings.Split(token, ".")
	bytes, e := jwt.DecodeSegment(payload[1])
	if e != nil {
		return user, e
	}
	content := ""
	for i := 0; i < len(bytes); i++ {
		content += string(bytes[i])
	}
	split := strings.Split(content, ",")
	id := strings.SplitAfter(split[2], ":")

	i := strings.Split(id[1], "\"")

	ID, err := strconv.Atoi(i[1])
	if err != nil {
		return user, err
	}
	return models.FindUserById(ID)
}
