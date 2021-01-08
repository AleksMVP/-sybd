package service

import (
	"github.com/AleksMVP/sybd/models"
)

type IServiceRepository interface {
	GetStatus() (status models.Status)
	Clear()
}