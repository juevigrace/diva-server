package jwt

import (
    "fmt"
    "os"
    "time"

    "github.com/golang-jwt/jwt/v5"
    "github.com/google/uuid"
    "github.com/juevigrace/diva-server/pkg/errs"
    "github.com/juevigrace/diva-server/pkg/config"
)

type JwtData struct {
    Token  *jwt.Token
    Claims *JWTClaims
}

type userClaims struct {
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
    jwtSecret   string           = os.Getenv("JWT_SECRET")
    Audience    jwt.ClaimStrings = jwt.ClaimStrings{"api"}
    AccessExp   = config.GetEnvOrDefault("JWT_ACCESS_TOKEN_EXP", 3600)
    RefreshExp  = config.GetEnvOrDefault("JWT_REFRESH_TOKEN_EXP", 86400)
)

func CreateToken(sessionId uuid.UUID, duration time.Duration) (string, error) {
    expiration := time.Now().UTC().Add(duration)
    claims := &userClaims{
        SessionID: sessionId,
    }
    return createJWT(claims, expiration)
}

func createJWT(claims *userClaims, expiration time.Time) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaims{
        *claims,
        jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expiration),
            IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
            NotBefore: jwt.NewNumericDate(time.Now().UTC()),
            Issuer:    Issuer,
            Subject:   claims.SessionID.String(),
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
        return nil, errs.ErrTokenNotValid
    }

    if len(claims.Audience) > 1 || claims.Audience[0] != "api" {
        return nil, errs.ErrBadAudience
    }

    if claims.Issuer != Issuer {
        return nil, errs.ErrBadIssuer
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
