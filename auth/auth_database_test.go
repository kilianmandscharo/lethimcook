package auth

import (
	"testing"

	"github.com/kilianmandscharo/lethimcook/errutil"
	"github.com/stretchr/testify/assert"
)

func TestCreateAdmin(t *testing.T) {
	// Given
	db := newTestAuthDatabase()

	// When
	admin := newTestAdmin()
	err := db.createAdmin(&admin)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, uint(1), admin.ID)

	// When
	err = db.createAdmin(&admin)

	// Then
	assert.ErrorIs(t, err, errutil.AuthErrorAdminAlreadyExists)
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
	admin := newTestAdmin()
	assert.NoError(t, db.createAdmin(&admin))

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
	assert.ErrorIs(t, err, errutil.AuthErrorNoAdminFound)

	// Given
	admin := newTestAdmin()
	assert.NoError(t, db.createAdmin(&admin))

	// When
	retrievedAdmin, err := db.readAdmin()

	// Then
	assert.NoError(t, err)
	assert.Equal(t, admin.ID, retrievedAdmin.ID)
	assert.Equal(t, admin.PasswordHash, retrievedAdmin.PasswordHash)
}

func TestReadAdminPasswordHash(t *testing.T) {
	// Given
	db := newTestAuthDatabase()

	// When
	_, err := db.readAdmin()

	// Then
	assert.ErrorIs(t, err, errutil.AuthErrorNoAdminFound)

	// Given
	admin := newTestAdmin()
	assert.NoError(t, db.createAdmin(&admin))

	// When
	passwordHash, err := db.readAdminPasswordHash()

	// Then
	assert.NoError(t, err)
	assert.Equal(t, admin.PasswordHash, passwordHash)
}

func TestUpdateAdminPassortHash(t *testing.T) {
	// Given
	db := newTestAuthDatabase()
	updatedHash := "updated test hash"

	// When
	err := db.updateAdminPasswordHash(updatedHash)

	// Then
	assert.ErrorIs(t, err, errutil.AuthErrorNoAdminFound)

	// Given
	admin := newTestAdmin()
	assert.NoError(t, db.createAdmin(&admin))

	// When
	err = db.updateAdminPasswordHash(updatedHash)

	// Then
	assert.NoError(t, err)

	// When
	retrievedAdmin, err := db.readAdmin()

	// Then
	assert.NoError(t, err)
	assert.Equal(t, updatedHash, retrievedAdmin.PasswordHash)
}
