package jwt_helper

import (
	"context"
	"errors"
	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/jwt"
	"os"
	user2 "otus_sn_wsserver/internal/models/user"
	"strconv"
	"strings"
	"time"
)

const (
	JWT_ERROR_KEY = "JWT_ERROR"
	JWT_USER_KEY  = "JWT_USER"
)

var TokenAuth *jwtauth.JWTAuth

type JWTToken struct {
	Token     string
	ExpiresIn time.Time
}

func Init() {
	secret, _ := os.LookupEnv("JWT_SECRET")
	TokenAuth = jwtauth.New("HS256", []byte(secret), nil)
}

func GenerateJwtToken(user *user2.User) *JWTToken {
	ttlStr, _ := os.LookupEnv("JWT_TOKEN_TTL")
	ttl, _ := strconv.Atoi(ttlStr)
	expiresIn := time.Now().Local().Add(time.Second * time.Duration(ttl))
	_, tokenString, _ := TokenAuth.Encode(map[string]interface{}{
		"user_id":    user.Id,
		"expires_in": expiresIn.UTC().Format(time.RFC3339),
	})
	return &JWTToken{
		Token:     tokenString,
		ExpiresIn: expiresIn,
	}
}

func GetCurrentUser(ctx context.Context) (user *user2.User, e error) {
	defer func() {
		if err := recover(); err != nil {
			user = nil
			e = err.(error)
		}
	}()
	user = ctx.Value(JWT_USER_KEY).(*user2.User)
	return user, nil
}

func GetTokenError(ctx context.Context) (errorString string, e error) {
	defer func() {
		if err := recover(); err != nil {
			errorString = ""
			e = err.(error)
		}
	}()
	errorString = ctx.Value(JWT_ERROR_KEY).(string)
	return errorString, nil
}

func GetTokenFromBearer(bearer string) string {
	if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
		return bearer[7:]
	}
	return ""
}

func GetUserIdFromToken(tokenString string) (int, error) {
	if tokenString != "" {
		token, err := jwtauth.VerifyToken(TokenAuth, tokenString)
		if err == nil {
			claims, err2 := token.AsMap(context.Background())
			if err2 == nil && jwt.Validate(token) == nil {
				expiresIn, err3 := time.Parse(time.RFC3339, claims["expires_in"].(string))
				if err3 == nil {
					if time.Now().UTC().After(expiresIn) {
						return 0, errors.New("Token expired")
					} else {
						return int(claims["user_id"].(float64)), nil
					}
				} else {
					return 0, err3
				}
			} else {
				return 0, err2
			}
		}
		return 0, err
	}
	return 0, errors.New("Token is empty")
}
