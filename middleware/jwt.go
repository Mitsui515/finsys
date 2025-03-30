package middleware

import (
	"context"
	"time"

	"github.com/Mitsui515/finsys/model"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/golang-jwt/jwt/v5"
)

type JWTConfig struct {
	SecretKey     string
	TokenExpire   time.Duration
	TokenIssuer   string
	TokenSubject  string
	TokenAudience string
}

func DefaultJWTConfig() *JWTConfig {
	return &JWTConfig{
		SecretKey:     "HKU_Project",
		TokenExpire:   24 * 7 * time.Hour,
		TokenIssuer:   "finsys",
		TokenSubject:  "auth",
		TokenAudience: "user",
	}
}

type CustomClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateToken(user *model.User) (string, error) {
	config := DefaultJWTConfig()
	claims := CustomClaims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.TokenExpire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    config.TokenIssuer,
			Subject:   config.TokenIssuer,
			Audience:  jwt.ClaimStrings{config.TokenAudience},
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.SecretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ParseToken(tokenString string) (*CustomClaims, error) {
	config := DefaultJWTConfig()
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.SecretKey), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrSignatureInvalid
}

func JWTAuth() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		tokenString := c.Request.Header.Get("Authorization")
		if tokenString == "" {
			c.JSON(consts.StatusUnauthorized, utils.H{
				"code":    consts.StatusUnauthorized,
				"message": "Unauthorized",
				"details": "Missing authorization token",
			})
			c.Abort()
			return
		}
		if len(tokenString) > 0 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}
		claims, err := ParseToken(tokenString)
		if err != nil {
			c.JSON(consts.StatusForbidden, utils.H{
				"code":    consts.StatusForbidden,
				"message": "Forbidden",
				"details": "Invalid or expired token",
			})
			c.Abort()
			return
		}
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Next(ctx)
	}
}
