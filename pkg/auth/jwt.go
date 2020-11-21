package auth

import (
	"github.com/spf13/viper"
	"github.com/veith/fgs-lib/pkg/ulid"
	"github.com/veith/goJWT/crypto"
	"github.com/veith/goJWT/jws"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"log"
	"time"
)

// user profile for internal use only
type Profile struct {
	Username string   `json:"username,omitempty"`
	Pwhash   string   `json:"pwhash,omitempty"`
	Roles    []string `json:"roles,omitempty"`
}

// Passwort hash erzeugen
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 5)
	return string(bytes), err
}

// Passwort hash prüfen
// true => passwort io, false => stimmt nicht
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// create a jwt token based on user profile
func CreateJWT(user *Profile) string {
	bytes, err := ioutil.ReadFile(viper.GetString("keys.private"))
	if err != nil {
		log.Fatal("Error reading private key")

	}
	rsaPrivate, keyErr := crypto.ParseRSAPrivateKeyFromPEM(bytes)
	if keyErr != nil {
		log.Fatal("Error parsing private key")
	}

	claims := jws.Claims{}
	claims.Set("rol", "user")
	claims.SetIssuer(viper.GetString("server.jwt.issuer"))
	claims.SetSubject(user.Username)
	now := time.Now()
	claims.SetIssuedAt(now)
	claims.SetExpiration(now.Add(time.Second * viper.GetDuration("server.jwt.expiry_duration_in_s")))
	claims.SetNotBefore(now)
	claims.SetJWTID(ulid.GenerateStringULID())
	jwt := jws.NewJWT(claims, crypto.SigningMethodRS256)
	byteToken, err := jwt.Serialize(rsaPrivate)

	if err != nil {
		log.Fatal("Error signing the key. ", err)
	}

	return string(byteToken)
}
