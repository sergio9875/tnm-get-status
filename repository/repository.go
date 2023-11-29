package repository

import "malawi-getstatus/models"

// Repository represent the repositories
type Repository interface {
	Close() error
	UpdateTransRefund(transId int, amount float64, activeStatus int, transrId int) error
	GetTransStatus(transId int) (*models.TransEntity, error)
}
