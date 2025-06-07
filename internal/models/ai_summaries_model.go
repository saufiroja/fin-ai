package models

import (
	"time"

	"github.com/saufiroja/fin-ai/internal/constants"
)

type AISummary struct {
	SummaryId   string               `json:"summary_id"`
	UserId      string               `json:"user_id"`
	PeriodType  constants.PeriodType `json:"period_type"`
	PeriodStart time.Time            `json:"period_start"`
	PeriodEnd   time.Time            `json:"period_end"`
	SummaryData any                  `json:"summary_data"` // type data jsonb for summary data
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
}
