package helpers

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/hasanbakirci/doc-system/internal/config"
	"github.com/hasanbakirci/doc-system/pkg/response"
	"github.com/labstack/gommon/log"
	"net/http"
	"time"
)

type UserClaim struct {
	ID   string `json:"id"`
	Role string `json:"role"`
	jwt.StandardClaims
}

func GenerateJwtToken(id, role string, cfg config.JwtSettings) string {
	secretKey := []byte(cfg.SecretKey)

	claims := UserClaim{
		ID:   id,
		Role: role,
		StandardClaims: jwt.StandardClaims{
			Audience:  "hasan@hasan.com",
			ExpiresAt: time.Now().Add(time.Duration(cfg.SessionTime) * time.Hour).Unix(),
			Issuer:    "hasan@hasan.com",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		log.Error("Couldn't get signed token : %v", err)
		response.Panic(http.StatusBadRequest, fmt.Sprintf("Couldn't get signed token : %v", err))
	}

	log.Info("The token was successfully generated.")
	return tokenString
}

func VerifyToken(token string, secret string) *UserClaim {
	secretKey := []byte(secret)

	decodedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			//return nil, fmt.Errorf("There was an error in parsing")
			response.Panic(http.StatusBadRequest, "There was an error in parsing")
		}
		return secretKey, nil
	})

	if err != nil {
		log.Error("Jwt token parse error : %v", err)
		response.Panic(http.StatusUnauthorized, fmt.Sprintf("Jwt token parse error : %v", err))
	}
	if !decodedToken.Valid {
		log.Error("Jwt token not valid ")
		response.Panic(http.StatusUnauthorized, "Jwt token not valid")
	}

	claims, ok := decodedToken.Claims.(jwt.MapClaims)
	if !ok {
		log.Error("Claims not found")
		response.Panic(http.StatusUnauthorized, "Claims not found")
	}

	userClaims := new(UserClaim)
	jsonString, _ := json.Marshal(claims)
	_ = json.Unmarshal(jsonString, &userClaims)

	log.Info("Successfully generated UserClaims from the token.")
	return userClaims
}
