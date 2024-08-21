package routes

import (
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

func generateToken(username string, tokenExpiryTime int32, jwtSecret string) (string, int64, error) {
	expiryDuration := 24 * time.Hour
	if tokenExpiryTime > 0 {
		expiryDuration = time.Duration(tokenExpiryTime) * time.Second
	}

	expiryTime := jwt.NewNumericDate(time.Now().Add(expiryDuration))
	expiryEpoch := expiryTime.Time.Unix()

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "jpleatherland/spacetraders",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: expiryTime,
		Subject:   username,
	})

	signedToken, err := newToken.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", 0, err
	}
	return signedToken, expiryEpoch, nil
}

func createSessionCookie(token string, expiryTime int64) http.Cookie {
	cookie := http.Cookie{
		Name:     "spacetradersSession",
		Value:    token,
		Expires:  time.Unix(expiryTime, 0),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	return cookie
}
