package clickhouse

import (
	"api/internal/entity"
)

// Migrations is a list of all gorm migrations for the clickhouse database.
var Migrations = []interface{}{
	&entity.Click{},
	&entity.Impression{},
}
