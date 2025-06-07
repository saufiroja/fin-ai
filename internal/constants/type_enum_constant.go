package constants

type TypeCategory string

const (
	IncomeCategory  TypeCategory = "income"
	ExpenseCategory TypeCategory = "expense"
)

type PeriodType string

const (
	PeriodTypeDaily   PeriodType = "daily"
	PeriodTypeWeekly  PeriodType = "weekly"
	PeriodTypeMonthly PeriodType = "monthly"
	PeriodTypeYearly  PeriodType = "yearly"
)

type RecommendationPriority string

const (
	RecommendationPriorityLow    RecommendationPriority = "low"
	RecommendationPriorityMedium RecommendationPriority = "medium"
	RecommendationPriorityHigh   RecommendationPriority = "high"
)

type RecommendationType string

const (
	RecommendationTypeBudgetAlert     RecommendationType = "budget alert"
	RecommendationTypeSavingTips      RecommendationType = "saving tip"
	RecommendationTypeSpendingWarning RecommendationType = "spending warning"
)
