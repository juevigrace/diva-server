package util

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JwtData struct {
	Token  *jwt.Token
	Claims *JWTClaims
}

type userClaims struct {
	UserId    uuid.UUID `json:"user_id"`
	SessionID uuid.UUID `json:"session_id"`
}

type JWTClaims struct {
	userClaims
	jwt.RegisteredClaims
}

const (
	Issuer string = "DIVA"
)

var (
	jwtSecret string           = os.Getenv("JWT_SECRET")
	Audience  jwt.ClaimStrings = jwt.ClaimStrings{
		"api",
	}
)

func CreateAccessToken(userId, sessionId uuid.UUID) (string, error) {
	var accessExpiration time.Time = time.Now().UTC().Add(1 * time.Hour)
	claims := &userClaims{
		UserId:    userId,
		SessionID: sessionId,
	}
	return createJWT(claims, accessExpiration)
}

func CreateRefreshToken(userId, sessionId uuid.UUID) (string, error) {
	var refreshExpiration time.Time = time.Now().UTC().Add(24 * time.Hour)
	claims := &userClaims{
		UserId:    userId,
		SessionID: sessionId,
	}
	return createJWT(claims, refreshExpiration)
}

func CreateResetToken(userId, sessionId uuid.UUID) (string, error) {
	var resetExpiration time.Time = time.Now().UTC().Add(10 * time.Minute)
	claims := &userClaims{
		UserId:    userId,
		SessionID: sessionId,
	}
	return createJWT(claims, resetExpiration)
}

func createJWT(claims *userClaims, expiration time.Time) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaims{
		*claims,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			NotBefore: jwt.NewNumericDate(time.Now().UTC()),
			Issuer:    Issuer,
			Subject:   claims.UserId.String(),
			Audience:  Audience,
		},
	})
	return token.SignedString([]byte(jwtSecret))
}

func ValidateJWT(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid reset token")
	}

	if len(claims.Audience) > 1 || claims.Audience[0] != "api" {
		return nil, errors.New("permission denied")
	}

	if claims.Issuer != Issuer {
		return nil, errors.New("permission denied")
	}

	return claims, nil
}

func GetTokenExpiration(tokenString string) (*time.Time, error) {
	claims, err := ValidateJWT(tokenString)
	if err != nil {
		return nil, err
	}
	return &claims.ExpiresAt.Time, nil
}
