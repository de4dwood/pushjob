package push

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

type Label struct {
	Key   string
	Value string
}

type PushJobStatus struct {
	labels         []Label
	StatusCode     int
	PushGatewayUrl string
	StartTime      time.Time
	EndTime        time.Time
	Duration       time.Duration
	JobName        string
}

var (
	jobstatus = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "push_job_status",
		Help: "status of the job ran in push infra",
	})
	completionTime = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "push_job_completionTime",
		Help: "The timestamp of the last completion of job, successful or not.",
	})
	successTime = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "push_job_successTime",
		Help: "The timestamp of the last successful job",
	})
	duration = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "push_job_duration",
		Help: "The duration of the last job in seconds.",
	})
)

func (p *PushJobStatus) Push() error {
	registry := prometheus.NewRegistry()
	registry.MustRegister(jobstatus, completionTime, duration)
	pusher := push.New(p.PushGatewayUrl, p.JobName).Gatherer(registry)
	jobstatus.Set(float64(p.StatusCode))
	completionTime.Set(float64(p.EndTime.Unix()))
	successTime.Set(float64(p.EndTime.Unix()))
	duration.Set(float64(p.EndTime.Sub(p.StartTime)))
	if p.StatusCode == 0 {
		pusher.Collector(successTime)
	}
	for _, s := range p.labels {
		pusher.Grouping(s.Key, s.Value)
	}
	if err := pusher.Add(); err != nil {
		return fmt.Errorf("could not push to pushgateway: %v", err)
	}
	return nil
}

func (p *PushJobStatus) AddLabel(l Label) {
	p.labels = append(p.labels, l)
}

func (p *PushJobStatus) GetLabels() string {
	s := ""
	for _, i := range p.labels {
		s += i.Key + "=" + i.Value + ","
	}
	return s
}
