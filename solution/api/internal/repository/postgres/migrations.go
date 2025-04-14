package postgres

import (
	"api/internal/entity"
)

// Migrations is a list of all gorm migrations for the postgres database.
var Migrations = []interface{}{
	&entity.Client{},
	&entity.Advertiser{},
	&entity.MLScore{},
	&entity.Campaign{},
	&entity.Targeting{},
}
