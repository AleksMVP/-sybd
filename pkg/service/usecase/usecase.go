package usecase

import (
	"github.com/AleksMVP/sybd/models"
	"github.com/AleksMVP/sybd/pkg/service"
)

type ServiceUseCase struct {
	serviceRepository service.IServiceRepository
}

func NewServiceUseCase(serviceRepository service.IServiceRepository) ServiceUseCase {
	return ServiceUseCase{
		serviceRepository: serviceRepository,
	}
}

func (h *ServiceUseCase)GetStatus() (status models.Status) {
	return h.serviceRepository.GetStatus()
}

func (h *ServiceUseCase)Clear() {
	h.serviceRepository.Clear()
}