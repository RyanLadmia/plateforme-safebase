package security

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Structure des données contenues dans le token
type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// Génère un token JWT signé avec la clé secrète
func GenerateJWT(secret string, userID uint, email, role string, duration time.Duration) (string, error) {
	// Crée les claims (données embarquées dans le token)
	claims := &Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)), // durée de validité (qui sera définie dans le user_service.go)
			IssuedAt:  jwt.NewNumericDate(time.Now()),               // date de création
			Issuer:    "safebase",                                   // nom du service
		},
	}

	// Crée le token signé
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// Vérifie et décode un token JWT
func VerifyJWT(secret, tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	// Si tout est ok, on retourne les données du token
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}
