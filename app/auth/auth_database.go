package auth

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/kilianmandscharo/lethimcook/errutil"
	"github.com/kilianmandscharo/lethimcook/logging"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type admin struct {
	ID           uint
	PasswordHash string
}

func newAdmin(passwordHash string) admin {
	return admin{PasswordHash: passwordHash}
}

type authDatabase struct {
	handler *gorm.DB
	logger  *logging.Logger
}

func NewAuthDatabase(logger *logging.Logger) *authDatabase {
	db, err := gorm.Open(sqlite.Open("./auth.db"), &gorm.Config{})
	if err != nil {
		logger.Fatal("failed to connect auth database: ", err)
	}
	db.AutoMigrate(&admin{})
	return &authDatabase{handler: db, logger: logger}
}

func (db *authDatabase) createAdmin(admin *admin) error {
	doesAdminExist, err := db.doesAdminExist()
	if err != nil {
		return errutil.AddMessageToAppError(err, "failed at createAdmin()")
	}
	if doesAdminExist {
		return &errutil.AppError{
			UserMessage: "Ein Admin existiert bereits",
			Err:         errors.New("failed at createAdmin(), an admin already exists"),
			StatusCode:  http.StatusConflict,
		}
	}
	if err := db.handler.Create(admin).Error; err != nil {
		return &errutil.AppError{
			UserMessage: "Datenbankfehler",
			Err: fmt.Errorf(
				"failed at createAdmin() with passwordHash, database failure: %w",
				err,
			),
			StatusCode: http.StatusInternalServerError,
		}
	}
	return nil
}

func (db *authDatabase) doesAdminExist() (bool, error) {
	if err := db.handler.First(&admin{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, &errutil.AppError{
			UserMessage: "Datenbankfehler",
			Err: fmt.Errorf(
				"failed at doesAdminExist(), database failure: %w",
				err,
			),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return true, nil
}

func (db *authDatabase) readAdmin() (admin, error) {
	var admin admin
	if err := db.handler.First(&admin).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return admin, &errutil.AppError{
				UserMessage: "Kein Admin gefunden",
				Err:         errors.New("failed at readAdmin(), no admin found"),
				StatusCode:  http.StatusNotFound,
			}
		}
		return admin, &errutil.AppError{
			UserMessage: "Datenbankfehler",
			Err: fmt.Errorf(
				"failed at readAdmin(), database failure: %w",
				err,
			),
			StatusCode: http.StatusInternalServerError,
		}
	}
	return admin, nil
}

func (db *authDatabase) readAdminPasswordHash() (string, error) {
	admin, err := db.readAdmin()
	if err != nil {
		return "", errutil.AddMessageToAppError(err, "failed at readAdminPasswordHash()")
	}
	return admin.PasswordHash, nil
}

func (db *authDatabase) updateAdminPasswordHash(newPasswordHash string) error {
	admin, err := db.readAdmin()
	if err != nil {
		return errutil.AddMessageToAppError(err, "failed at updateAdminPasswordHash()")
	}
	admin.PasswordHash = newPasswordHash
	if err := db.handler.Save(&admin).Error; err != nil {
		return &errutil.AppError{
			UserMessage: "Datenbankfehler",
			Err: fmt.Errorf(
				"failed at updateAdminPasswordHash() with newPasswordHash, database failure: %w",
				err,
			),
			StatusCode: http.StatusInternalServerError,
		}
	}
	return nil
}
