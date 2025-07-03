package seed

import (
	"super-web-server/internal/model"
	"super-web-server/pkg/database"
)

func SeedUserRole(db *database.DB) error {
	tx := db.Begin()
	for _, role := range model.UserRoleSet {
		var count int64
		tx.Model(&model.UserRole{}).Where("code = ?", role.Code).Count(&count)
		if count > 0 {
			continue
		}
		if err := tx.Create(&role).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
