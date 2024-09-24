package auth

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/kilianmandscharo/lethimcook/errutil"
	"github.com/kilianmandscharo/lethimcook/logging"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func newTestAuthDatabase() authDatabase {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		fmt.Println("failed to connect test database: ", err)
		os.Exit(1)
	}
	db.Migrator().DropTable(&admin{})
	db.AutoMigrate(&admin{})
	logger := logging.New(logging.Debug)
	return authDatabase{handler: db, logger: &logger}
}

func newTestAdmin() admin {
	return admin{PasswordHash: "test hash"}
}

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
	assert.Error(t, err)
	appError, ok := err.(*errutil.AppError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusConflict, appError.StatusCode)
	assert.Equal(t, "Ein Admin existiert bereits", appError.UserMessage)
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
	assert.Error(t, err)
	appError, ok := err.(*errutil.AppError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusNotFound, appError.StatusCode)
	assert.Equal(t, "Kein Admin gefunden", appError.UserMessage)

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
	assert.Error(t, err)
	appError, ok := err.(*errutil.AppError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusNotFound, appError.StatusCode)
	assert.Equal(t, "Kein Admin gefunden", appError.UserMessage)

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
	assert.Error(t, err)
	appError, ok := err.(*errutil.AppError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusNotFound, appError.StatusCode)
	assert.Equal(t, "Kein Admin gefunden", appError.UserMessage)

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
