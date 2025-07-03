package seed

import (
	"super-web-server/pkg/database"
	"super-web-server/pkg/snowflake"
)

func Seed(db *database.DB, snowflake *snowflake.Snowflake) error {
	if err := SeedUserRole(db); err != nil {
		return err
	}

	if err := SeedUser(db, snowflake); err != nil {
		return err
	}

	return nil
}
