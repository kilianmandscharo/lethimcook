package auth

import (
	"errors"
	"log"

	"github.com/kilianmandscharo/lethimcook/errutil"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
}

func newTestAuthDatabase() authDatabase {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		log.Fatal("failed to connect test database")
	}

	db.Migrator().DropTable(&admin{})

	db.AutoMigrate(&admin{})

	return authDatabase{handler: db}
}

func newAuthDatabase() authDatabase {
	db, err := gorm.Open(sqlite.Open("./auth.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		log.Fatal("failed to connect auth database")
	}

	db.AutoMigrate(&admin{})

	return authDatabase{handler: db}
}

func (db *authDatabase) doesAdminExist() (bool, error) {
	if err := db.handler.First(&admin{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (db *authDatabase) createAdmin(admin *admin) error {
	doesAdminExist, err := db.doesAdminExist()
	if err != nil {
		return err
	}
	if doesAdminExist {
		return errutil.AuthErrorAdminAlreadyExists
	}

	if err := db.handler.Create(admin).Error; err != nil {
		return err
	}

	return nil
}

func (db *authDatabase) readAdmin() (admin, errutil.AuthError) {
	var admin admin

	if err := db.handler.First(&admin).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return admin, errutil.AuthErrorNoAdminFound
		}
		return admin, errutil.AuthErrorDatabaseFailure
	}

	return admin, nil
}

func (db *authDatabase) readAdminPasswordHash() (string, errutil.AuthError) {
	admin, err := db.readAdmin()

	if err != nil {
		return "", err
	}

	return admin.PasswordHash, nil
}

func (db *authDatabase) updateAdminPasswordHash(newPasswordHash string) errutil.AuthError {
	admin, err := db.readAdmin()
	if err != nil {
		return err
	}

	admin.PasswordHash = newPasswordHash
	if err := db.handler.Save(&admin).Error; err != nil {
		return errutil.AuthErrorDatabaseFailure
	}

	return nil
}
