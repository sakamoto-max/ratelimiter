package jwt

import (
	"fmt"
	"time"

	"github.com/sakamoto-max/ratelimiter/internal/config"
	myErr "github.com/sakamoto-max/ratelimiter/internal/pkg/myerrors"

	"github.com/golang-jwt/jwt/v5"
)

var SECRETKEY string

type Claims struct {
	Ownername string `json:"ownerName"`
	jwt.RegisteredClaims
}

var (
	DefaultExpiresAt = time.Now().Add(time.Hour * 8760 * 100)
)

func GenerateToken(ownername string, expiresAt time.Time) (string, error) {

	claims := &Claims{
		Ownername: ownername,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "rate_limiter",
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString([]byte(SECRETKEY))
	if err != nil {
		return "", myErr.WrapErr(fmt.Errorf("failed to generate token : %w", err), myErr.InternalServerErr)
	}

	return tokenStr, nil
}

func ValidateToken(userSentToken string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(userSentToken, claims, func(t *jwt.Token) (any, error) {
		return []byte(SECRETKEY), nil
	})

	if err != nil {
		return nil, myErr.WrapErr(fmt.Errorf("failed to parse token : %w", err), myErr.UnauthorizedErr)
	}

	if !token.Valid {
		return nil, myErr.WrapErr(fmt.Errorf("token is invalid"), myErr.UnauthorizedErr)
	}

	return claims, nil
}

func Init(config *config.Config) {
	SECRETKEY = config.Auth.SecretKey
}
