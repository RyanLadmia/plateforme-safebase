package security

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Structure of the data contained in the token
type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// Generate a signed JWT token with the secret key
func GenerateJWT(secret string, userID uint, email, role string, duration time.Duration) (string, error) {
	// Create the claims (data embedded in the token)
	claims := &Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)), // validity duration (which will be defined in user_service.go)
			IssuedAt:  jwt.NewNumericDate(time.Now()),               // creation date
			Issuer:    "safebase",                                   // service name
		},
	}

	// Create the signed token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// Verify and decode a JWT token
func VerifyJWT(secret, tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	// If everything is ok, return the token data
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}
