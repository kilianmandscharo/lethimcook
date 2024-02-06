package auth

// import (
// 	"testing"
//
// 	"github.com/stretchr/testify/assert"
// )
//
// func TestHashPassword(t *testing.T) {
// 	// Given
// 	authorizer := newTestAuthorizer()
//
// 	// When
// 	_, err := authorizer.hashPassword("test password")
//
// 	// Then
// 	assert.NoError(t, err)
// }
//
// func TestCheckPasswordHash(t *testing.T) {
// 	// Given
// 	authorizer := newTestAuthorizer()
//
// 	// When
// 	password := "test password"
// 	hash, _ := authorizer.hashPassword(password)
//
// 	// Then
// 	assert.True(t, authorizer.checkPasswordHash(password, hash))
// 	assert.False(t, authorizer.checkPasswordHash("invalid password", hash))
// }
