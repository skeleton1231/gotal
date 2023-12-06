package user

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/skeleton1231/gotal/internal/apiserver/store/model"
	"github.com/spf13/viper"
)

func generateJWTToken(user *model.User) (string, time.Time, error) {
	// 构造 JWT claims
	claims := jwt.MapClaims{
		"userID":   user.ID, // 或者使用其他用户唯一标识
		"exp":      time.Now().Add(viper.GetDuration("jwt.timeout")).Unix(),
		"identity": user.Name, // 根据 PayloadFunc 中的逻辑来设置
		// 可以添加更多的 claims
	}

	// 创建 JWT 令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用您的密钥签署令牌
	signedToken, err := token.SignedString([]byte(viper.GetString("jwt.key")))
	if err != nil {
		return "", time.Time{}, err
	}

	// 获取令牌的过期时间
	expireTime := time.Now().Add(viper.GetDuration("jwt.timeout"))

	return signedToken, expireTime, nil
}
