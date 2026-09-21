package services

import (
	"QuickGin/db"

	"QuickGin/models/cache"

	"github.com/jmoiron/sqlx"
)

type AdminServiceConfig struct {
	DB    *sqlx.DB
	cache cache.Cache
}
type AdminService struct {
	cfg AdminServiceConfig
}

func NewAdminService() *AdminService {
	return &AdminService{cfg: AdminServiceConfig{
		DB:    db.AppDB(),
		cache: db.AppCache(),
	}}
}
