package repository

// Repository represent the repositories
type Repository interface {
	Close() error
	UpdateTransRefund(transId int, amount float64, activeStatus int, transrId int) error
}
