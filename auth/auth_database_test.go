package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateAdmin(t *testing.T) {
	// Given
	db := newTestAuthDatabase()

	// When
	admin := Admin{PasswordHash: "test hash"}
	err := db.createAdmin(&admin)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, uint(1), admin.ID)

	// When
	admin = Admin{PasswordHash: "test hash"}
	err = db.createAdmin(&admin)

	// Then
	assert.Error(t, err)
}

func TestDoesAdminExist(t *testing.T) {
	// Given
	db := newTestAuthDatabase()

	// When
	doesAdminExist, err := db.doesAdminExist()

	// Then
	assert.NoError(t, err)
	assert.False(t, doesAdminExist)

	// Given
	admin := Admin{PasswordHash: "test hash"}
	db.createAdmin(&admin)

	// When
	doesAdminExist, err = db.doesAdminExist()

	// Then
	assert.NoError(t, err)
	assert.True(t, doesAdminExist)
}

func TestReadAdmin(t *testing.T) {
	// Given
	db := newTestAuthDatabase()

	// When
	_, err := db.readAdmin()

	// Then
	assert.Error(t, err)

	// Given
	admin := Admin{PasswordHash: "test hash"}
	db.createAdmin(&admin)

	// When
	retrievedAdmin, err := db.readAdmin()

	// Then
	assert.NoError(t, err)
	assert.Equal(t, admin.PasswordHash, retrievedAdmin.PasswordHash)
}

func TestReadAdminPasswordHash(t *testing.T) {
	// Given
	db := newTestAuthDatabase()

	// When
	_, err := db.readAdmin()

	// Then
	assert.Error(t, err)

	// Given
	admin := Admin{PasswordHash: "test hash"}
	db.createAdmin(&admin)

	// When
	passwordHash, err := db.readAdminPasswordHash()

	// Then
	assert.NoError(t, err)
	assert.Equal(t, admin.PasswordHash, passwordHash)
}

func TestUpdateAdminPassortHash(t *testing.T) {
	// Given
	db := newTestAuthDatabase()

	// When
	updatedHash := "updated test hash"
	err := db.updateAdminPasswordHash(updatedHash)

	// Then
	assert.Error(t, err)

	// Given
	admin := Admin{PasswordHash: "test hash"}
	db.createAdmin(&admin)

	// When
	err = db.updateAdminPasswordHash(updatedHash)
	retrievedAdmin, _ := db.readAdmin()

	// Then
	assert.NoError(t, err)
	assert.Equal(t, updatedHash, retrievedAdmin.PasswordHash)
}
