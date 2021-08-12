/*
 * @Author: your name
 * @Date: 2021-06-07 09:13:06
 * @LastEditTime: 2021-08-12 17:13:29
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /nethopper/utils/token.go
 */
package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

//CreateToken token create
func CreateToken(secret string, content string, expire time.Duration) (string, error) {
	claims := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(expire).Unix(),
		Subject:   content,
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return jwtToken.SignedString([]byte(secret))
}

//ValidToken token valid
func ValidToken(secret string, token string) (string, error) {
	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil || jwtToken.Claims.Valid() != nil {
		return "", err
	}
	claims := jwtToken.Claims.(jwt.MapClaims)
	if !claims.VerifyExpiresAt(time.Now().Unix(), true) {
		return "", errors.New("token expired")
	}
	return claims["sub"].(string), nil
}
