package monitor

import (
	"github.com/OlegVankov/verbose-umbrella/internal/logger"
	"github.com/OlegVankov/verbose-umbrella/internal/storage"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"time"
)

func SendMetrics(m *Monitor, addr string, reportInterval int) {
	client := resty.New()
	url := "http://" + addr + "/update"
	for {
		<-time.After(time.Duration(reportInterval) * time.Second)

		for _, body := range m.GetBody() {
			var metric storage.Metrics

			resp, err := client.R().
				SetHeader("Content-Type", "application/json").
				SetBody(body).
				SetResult(&metric).
				Post(url)

			if err != nil {
				logger.Log.Error("resty request error", zap.Error(err))
				break
			}

			logger.Log.Info("SendMetric", zap.String("URL", resp.Request.URL),
				zap.String("body", resp.String()),
				zap.String("StatusCode", resp.Status()),
				zap.Any("metric", metric))
		}
		// обнулим pollCounter после отправки метрик
		m.resetPollCount()
	}
}
