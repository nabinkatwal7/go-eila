package model

// ... existing ...

type AnomalySeverity string

const (
	SeverityHigh   AnomalySeverity = "High"
	SeverityMedium AnomalySeverity = "Medium"
	SeverityLow    AnomalySeverity = "Low"
)

type Anomaly struct {
	ID          int64 // Virtual ID usually
	Type        string // "Large Transaction", "Budget Drift"
	Description string
	Severity    AnomalySeverity
	Date        string
}
