package postgres

import "bot/internal/entity"

// Migrations is a list of all gorm migrations for the postgres database.
var Migrations = []interface{}{
	&entity.User{},
}
