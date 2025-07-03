package jwt

import (
	"errors"
	"strings"
	"super-web-server/internal/ctx"
	"super-web-server/internal/exception"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Config struct {
	Secret string
	Expire time.Duration
	Issuer string
}

type JWTClaims struct {
	jwt.RegisteredClaims
	UserUniqueID int64 `json:"userUniqueId"`
	RefreshAt    int64 `json:"refreshAt"`
}

type JWT struct {
	Config Config
}

func NewJWT(config Config) *JWT {
	return &JWT{
		Config: config,
	}
}

func (j *JWT) JWT() gin.HandlerFunc {
	return func(gtx *gin.Context) {

		appCtx := ctx.NewAppCtx(gtx)

		token, err := j.GetTokenFromGinContext(gtx)
		if err != nil {
			appCtx.ToError(exception.ExceptionTokenNotFound)
			return
		}

		claims, err := j.ParseAndVerifyToken(token)
		if err != nil {
			appCtx.ToError(exception.ExceptionTokenExpired)
			return
		}

		// if token is near to expire, generate a new token
		if time.Now().Unix() > claims.RefreshAt {
			newToken, err := j.GenerateToken(claims.UserUniqueID)
			if err != nil {
				appCtx.ToError(exception.ExceptionTokenGenerateFailed)
				return
			}
			appCtx.Header("New-Token", newToken)
		}

		appCtx.SetUserUniqueID(claims.UserUniqueID)
		appCtx.Next()
	}
}

func (j *JWT) GetTokenFromGinContext(c *gin.Context) (string, error) {
	var token string
	keys := []string{"Authorization", "token"}
	for _, key := range keys {
		if value, exist := c.GetQuery(key); exist {
			token = strings.TrimPrefix(value, "Bearer ")
			break
		}
		if value := c.GetHeader(key); value != "" {
			token = strings.TrimPrefix(value, "Bearer ")
			break
		}
	}
	if token == "" {
		return "", errors.New("token is empty")
	}
	return token, nil
}

func (j *JWT) GetSecret() []byte {
	return []byte(j.Config.Secret)
}

func (j *JWT) GetExpire() time.Duration {
	return j.Config.Expire
}

func (j *JWT) GetIssuer() string {
	return j.Config.Issuer
}

func (j *JWT) ExpireAt() time.Time {
	return time.Now().Add(j.GetExpire())
}

func (j *JWT) RefreshAt() time.Time {
	return time.Now().Add(j.GetExpire() / 3 * 2)
}

func (j *JWT) GenerateClaims(userUniqueID int64) JWTClaims {
	return JWTClaims{
		UserUniqueID: userUniqueID,
		RefreshAt:    j.RefreshAt().Unix(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(j.ExpireAt()),
			Issuer:    j.GetIssuer(),
		},
	}
}

func (j *JWT) GenerateToken(userUniqueID int64) (string, error) {
	claims := j.GenerateClaims(userUniqueID)
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tokenClaims.SignedString(j.GetSecret())
}

func (j *JWT) ParseAndVerifyToken(token string) (*JWTClaims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &JWTClaims{}, func(token *jwt.Token) (any, error) {
		return j.GetSecret(), nil
	})
	if err != nil {
		return nil, err
	}
	if tokenClaims != nil {
		claims, ok := tokenClaims.Claims.(*JWTClaims)
		if ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
