package auth

import (
	"errors"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Admin struct {
	ID           uint
	PasswordHash string
}

type AuthDatabase struct {
	handler *gorm.DB
}

func newTestAuthDatabase() AuthDatabase {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		log.Fatal("failed to connect test database")
	}

	db.Migrator().DropTable(&Admin{})

	db.AutoMigrate(&Admin{})

	return AuthDatabase{handler: db}
}

func newAuthDatabase() (AuthDatabase, error) {
	db, err := gorm.Open(sqlite.Open("./auth.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		return AuthDatabase{}, err
	}

	db.AutoMigrate(&Admin{})

	return AuthDatabase{handler: db}, nil
}

func (db *AuthDatabase) doesAdminExist() (bool, error) {
	if err := db.handler.First(&Admin{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (db *AuthDatabase) createAdmin(admin *Admin) error {
	doesAdminExist, err := db.doesAdminExist()

	if err != nil {
		return err
	}

	if doesAdminExist {
		return errors.New("there can only be one admin")
	}

	if err := db.handler.Create(admin).Error; err != nil {
		return err
	}

	return nil
}

func (db *AuthDatabase) readAdmin() (Admin, error) {
	var admin Admin

	if err := db.handler.First(&admin).Error; err != nil {
		return admin, err
	}

	return admin, nil
}

func (db *AuthDatabase) readAdminPasswordHash() (string, error) {
	admin, err := db.readAdmin()

	if err != nil {
		return "", err
	}

	return admin.PasswordHash, nil
}

func (db *AuthDatabase) updateAdminPasswordHash(newPasswordHash string) error {
	admin, err := db.readAdmin()

	if err != nil {
		return err
	}

	admin.PasswordHash = newPasswordHash

	if err := db.handler.Save(&admin).Error; err != nil {
		return err
	}

	return nil
}
