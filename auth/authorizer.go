package auth

// import (
// 	"log"
// 	"os"
// 	"time"
//
// 	"github.com/golang-jwt/jwt"
// 	"github.com/joho/godotenv"
// 	"golang.org/x/crypto/bcrypt"
// )
//
// type Authorizer struct {
// 	db AuthDatabase
// }
//
// func newTestAuthorizer() Authorizer {
// 	db := newTestAuthDatabase()
//
// 	return Authorizer{db: db}
// }
//
// func NewAuthorizer() Authorizer {
// 	db, err := newAuthDatabase()
//
// 	if err != nil {
// 		log.Fatal("failed to create auth database", err)
// 	}
//
// 	return Authorizer{db: db}
// }
//
// func (a *Authorizer) DoesAdminExist() bool {
// 	doesAdminExist, err := a.db.doesAdminExist()
//
// 	if err != nil {
// 		log.Fatal("failed to read admin from database")
// 	}
//
// 	return doesAdminExist
// }
//
// func (a *Authorizer) CreateAdmin(password string) {
// 	if len([]byte(password)) > 72 {
// 		log.Fatal("the password length must be less than 72 bytes")
// 	}
//
// 	hash, err := a.hashPassword(password)
//
// 	if err != nil {
// 		log.Fatal("failed to hash password", err)
// 	}
//
// 	err = a.db.createAdmin(&Admin{PasswordHash: hash})
//
// 	if err != nil {
// 		log.Fatal("failed to create admin", err)
// 	}
// }
//
// func (a *Authorizer) login(password string) bool {
// 	hash, err := a.db.readAdminPasswordHash()
//
// 	if err != nil {
// 		log.Fatal("failed to read admin password hash", err)
// 	}
//
// 	return a.checkPasswordHash(password, hash)
// }
//
// func (a *Authorizer) createToken(key string) (string, error) {
// 	claims := &jwt.StandardClaims{
// 		ExpiresAt: time.Now().Add(60 * time.Minute).UnixMilli(),
// 	}
//
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
//
// 	tokenString, err := token.SignedString([]byte(key))
//
// 	if err != nil {
// 		return "", err
// 	}
//
// 	return tokenString, nil
// }
//
// func (a *Authorizer) hashPassword(password string) (string, error) {
// 	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 13)
// 	return string(bytes), err
// }
//
// func (a *Authorizer) checkPasswordHash(password, hash string) bool {
// 	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
// 	return err == nil
// }
//
// func (a *Authorizer) LoadPrivateKeyFromEnv() string {
// 	err := godotenv.Load(".env")
//
// 	if err != nil {
// 		log.Fatal("failed to load .env file", err)
// 	}
//
// 	return os.Getenv("JWT_PRIVATE_KEY")
// }
