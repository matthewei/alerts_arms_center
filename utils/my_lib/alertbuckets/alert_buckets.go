package alertbuckets

import (
	alertmanagertemplate "github.com/prometheus/alertmanager/template"
	"time"
)

// Alert is an alert received from the Alertmanager
type Alert = alertmanagertemplate.Alert

// Alerts is a set of alerts received from the Alertmanager
type Alerts = alertmanagertemplate.Alerts

// AlertsData is the alerts object received from the Alertmanager
type AlertsData = alertmanagertemplate.Data

// AlertBucket mutualizes alerts by Trap IDs
type AlertBuckets struct {
	AlertBuckets []AlertBucket
}

type AlertBucket struct {
	Status       string
	Annotations  map[string]string
	Labels       map[string]string
	StartsAt     time.Time
	EndsAt       time.Time
	GeneratorURL string
}

func New(alerts Alerts) *AlertBuckets {
	buckets := new(AlertBuckets)
	for _, alertDetails := range alerts {
		buckets.AlertBuckets = append(buckets.AlertBuckets, AlertBucket{
			Status:      alertDetails.Status,
			Labels:      alertDetails.Labels,
			Annotations: alertDetails.Annotations,
			StartsAt:    alertDetails.StartsAt,
			EndsAt:      alertDetails.EndsAt,
		})
	}
	return buckets
}
