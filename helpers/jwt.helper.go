
package helpers

import (
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	Username	string	`json:"username"`
	jwt.StandardClaims
}

func JWTGenerateTokenFromUser(user models.User, validUntil time.Duration) (string, error) {
	
}
 
func JWTParseToken(token string) (jwt.StandardClaims, error) {

}
