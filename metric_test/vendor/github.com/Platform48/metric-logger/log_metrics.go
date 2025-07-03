package metriclogger

import (
	"context"
	"fmt"
	"time"

	monitoring "cloud.google.com/go/monitoring/apiv3"
	"cloud.google.com/go/monitoring/apiv3/v2/monitoringpb"
	mpb "google.golang.org/genproto/googleapis/api/metric"
	gcprpb "google.golang.org/genproto/googleapis/api/monitoredres"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type MetricLogger struct {
	ctx       context.Context
	ProjectId string
}

type Metric struct {
	Name   string
	Labels map[string]string
}

func NewMetricLogger(ctx context.Context, projectId string) *MetricLogger {
	return &MetricLogger{ctx: ctx, ProjectId: projectId}
}

func (ml *MetricLogger) LogMetric(metric Metric) error {
	client, err := monitoring.NewMetricClient(ml.ctx)
	if err != nil {
		return fmt.Errorf("error occured when initilaising monitoring client: %v", err)
	}
	defer client.Close()

	now := timestamppb.New(time.Now())
	point := &monitoringpb.Point{
		Interval: &monitoringpb.TimeInterval{EndTime: now},
		Value: &monitoringpb.TypedValue{
			Value: &monitoringpb.TypedValue_Int64Value{
				Int64Value: 1,
			},
		},
	}

	ts := &monitoringpb.TimeSeries{
		Metric: &mpb.Metric{
			Type:   fmt.Sprintf("custom.googleapis.com/%s", metric.Name),
			Labels: metric.Labels,
		},
		Resource: &gcprpb.MonitoredResource{
			Type:   "global",
			Labels: map[string]string{"project_id": ml.ProjectId},
		},
		Points: []*monitoringpb.Point{point},
	}

	req := &monitoringpb.CreateTimeSeriesRequest{
		Name:       fmt.Sprintf("projects/%s", ml.ProjectId),
		TimeSeries: []*monitoringpb.TimeSeries{ts},
	}

	if err := client.CreateTimeSeries(ml.ctx, req); err != nil {
		return fmt.Errorf("metric logging error: %v", err)
	}

	return nil
}
