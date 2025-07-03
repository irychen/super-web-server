package seed

import (
	"super-web-server/internal/model"
	"super-web-server/pkg/database"
	"super-web-server/pkg/snowflake"
	"super-web-server/pkg/utils"
)

func SeedUser(db *database.DB, snowflake *snowflake.Snowflake) error {
	tx := db.Begin()

	const adminEmail = "admin@example.com"
	const adminPassword = "123456"

	adminRole := model.UserRole{}
	if err := tx.Model(&model.UserRole{}).Where("code = ?", model.UserRoleCodeSuperAdmin).First(&adminRole).Error; err != nil {
		tx.Rollback()
		return err
	}

	// check if admin user exists
	var count int64
	if err := tx.Model(&model.User{}).Where("email = ?", adminEmail).Count(&count).Error; err != nil {
		tx.Rollback()
		return err
	}
	if count > 0 {
		return nil
	}

	uniqueID := snowflake.GenerateID()

	salt, err := utils.GenerateSalt(6)
	if err != nil {
		tx.Rollback()
		return err
	}
	hashedPassword, err := utils.CryptHash(adminPassword, salt)
	if err != nil {
		tx.Rollback()
		return err
	}

	adminUser := model.User{
		UniqueID: uniqueID,
		Email:    adminEmail,
		Password: hashedPassword,
		Salt:     salt,
		Roles:    []*model.UserRole{&adminRole},
	}

	if err := tx.Create(&adminUser).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}
