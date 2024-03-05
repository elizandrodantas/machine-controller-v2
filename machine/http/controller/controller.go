package controller

import (
	"time"

	"github.com/elizandrodantas/machine-controller-v2/domain"
)

type MachineController struct {
	MachineUsecase domain.MachineUsecase
	RuleUsecase    domain.RuleUsecase
	Timeout        time.Duration
}
