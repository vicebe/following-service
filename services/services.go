package services

import (
	"log"

	"github.com/vicebe/following-service/data"
)

// AppService is register of all implemented services.
type AppService struct {
	*UserService
}

func NewAppService(l *log.Logger, ds *data.Store) *AppService {
	return &AppService{
		UserService: NewUserService(l, ds),
	}
}
