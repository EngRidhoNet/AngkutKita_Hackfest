// utils/jwt.go
package utils
import (
	"time"
	"github.com/dgrijalva/jwt-go"
)
var secretKey = []byte("your_secret_key") // Ganti dengan kunci rahasia yang kuat
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
func GenerateToken(username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}
