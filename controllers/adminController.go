package controllers

import (
	"QuickGin/services"
)

type AdminController struct {
	svc *services.AdminService
}

func NewAdminController() *AdminController {
	return &AdminController{
		svc: services.NewAdminService(),
	}
}
