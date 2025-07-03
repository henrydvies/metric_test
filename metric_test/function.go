package metrictester

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	metriclogger "github.com/Platform48/metric-logger"
)

func init() {
	functions.HTTP("MetricTest", MetricTest)
}

type Response struct {
	Greeting string `json:"greeting"`
}

func convertToStringMap(data map[string]interface{}) map[string]string {
	result := make(map[string]string)
	for key, value := range data {
		switch v := value.(type) {
		case string:
			result[key] = v
		case float64:
			result[key] = fmt.Sprintf("%v", v)
		case bool:
			result[key] = fmt.Sprintf("%v", v)
		case nil:
			result[key] = ""
		default:
			// Handle arrays, maps, etc. as needed
			result[key] = fmt.Sprintf("%v", v)
		}
	}
	return result
}

func MetricTest(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var jsonBody map[string]interface{}

	if err := json.NewDecoder(r.Body).Decode(&jsonBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	stringMap := convertToStringMap(jsonBody)

	metric_logging := metriclogger.NewMetricLogger(ctx, "p48-development")
	stats := metriclogger.Metric{Name: "test_stats", Labels: stringMap}

	if err := metric_logging.LogMetric(stats); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Errorf("Error occured when logging: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
} ///
