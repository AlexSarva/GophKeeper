package storage

import (
	"AlexSarva/GophKeeper/models"
	"errors"

	"github.com/google/uuid"
)

// ErrDuplicatePK error that occurs when adding exists user or order number
var ErrDuplicatePK = errors.New("duplicate PK")

// ErrNoValues error that occurs when no values selected from database
var ErrNoValues = errors.New("no values from select")

// Repo primary interface for all types of databases
type Repo interface {
	Ping() bool
	ServiceRegistration(serviceInfo models.InputService) error
	GetServiceInfo(serviceID uuid.UUID) (*models.Service, error)
	DeleteService(userID, serviceID uuid.UUID) error
	RecoveryService(userID, serviceID uuid.UUID) error
	GetServiceList(deleted bool) ([]*models.Service, error)
	EditService(userID uuid.UUID, serviceInfo models.InputService) error
	StartService(lunchedService models.LunchedService) error
	StopService(stoppedService models.StoppedService) error
}
