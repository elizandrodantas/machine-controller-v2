package controller

import (
	"time"

	"github.com/elizandrodantas/machine-controller-v2/domain"
)

type RuleController struct {
	RuleUsecase domain.RuleUsecase
	Timeout     time.Duration
}
