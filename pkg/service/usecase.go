package service

import (
	"github.com/AleksMVP/sybd/models"
)

type IServiceUseCase interface {
	GetStatus() (status models.Status)
	Clear()
}