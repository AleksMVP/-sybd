package delivery

import (
	"github.com/AleksMVP/sybd/pkg/service"
	"github.com/AleksMVP/sybd/utils"
	"net/http"
)

type ServiceDelivery struct {
	serviceUseCase service.IServiceUseCase
}

func NewServiceDelivery(serviceUseCase service.IServiceUseCase) ServiceDelivery {
	return ServiceDelivery{
		serviceUseCase: serviceUseCase,
	}
}

func (h *ServiceDelivery)PostServiceClear(w http.ResponseWriter, r *http.Request) {
	h.serviceUseCase.Clear()
	utils.WriteNotEasyJson(w, http.StatusOK, struct{}{})

	/*
	 * 200 -> Delete all user data
	 */
}

func (h *ServiceDelivery)GetServiceStatus(w http.ResponseWriter, r *http.Request) {
	status := h.serviceUseCase.GetStatus()
	utils.WriteJson(w, http.StatusOK, status)

	/*
	 * 200 -> Status
	 */
}