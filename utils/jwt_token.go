package utils

import (
  "encoding/base64"
  "fmt"
  "time"

  "github.com/golang-jwt/jwt"
  "github.com/google/uuid"
)

func ValidateToken(token string, publicKey string) (
  userId interface{},
  jti interface{},
  err error) {

  decodedPublicKey, err := base64.StdEncoding.DecodeString(publicKey)
  if err != nil {
    return
  }

  key, err := jwt.ParseRSAPublicKeyFromPEM(decodedPublicKey)
  if err != nil {
    return
  }

  parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
    if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
      return nil, fmt.Errorf("unexpected method: %s", t.Header["alg"])
    }
    return key, nil
  })

  if err != nil {
    return
  }

  claims, ok := parsedToken.Claims.(jwt.MapClaims)
  if !ok || !parsedToken.Valid {
    return userId, jti, fmt.Errorf("validate: invalid token")
  }

  userId = claims["sub"]
  jti = claims["jti"]

  return
}

func GenerateToken(ttlSec int, userId interface{}, privateKey string) (
  token string,
  jti string,
  expiredAt time.Time,
  err error) {

  decodedPrivateKey, err := base64.StdEncoding.DecodeString(privateKey)
  if err != nil {
    return
  }

  key, err := jwt.ParseRSAPrivateKeyFromPEM(decodedPrivateKey)
  if err != nil {
    return
  }

  now := time.Now().UTC()
  jti = uuid.New().String()
  expiredAt = now.Add(time.Second * time.Duration(ttlSec))

  claims := make(jwt.MapClaims)
  claims["sub"] = userId
  claims["exp"] = expiredAt.Unix()
  claims["iat"] = now.Unix()
  claims["nbf"] = now.Unix()
  claims["jti"] = jti

  token, err = jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)

  return
}
