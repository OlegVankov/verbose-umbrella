package monitor

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"

	"github.com/OlegVankov/verbose-umbrella/internal/logger"
	"github.com/OlegVankov/verbose-umbrella/internal/storage"
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
				SetHeader("Accept-Encoding", "gzip").
				SetHeader("Content-Encoding", "gzip").
				SetBody(body).
				SetResult(&metric).
				Post(url)

			if err != nil {
				logger.Log.Error("resty request error", zap.Error(err))
				continue
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

func SendBatch(m *Monitor, addr string, reportInterval int) {
	wait := []time.Duration{1, 3, 5}
	step := 0
	client := resty.New()
	client.SetRetryCount(3).SetRetryWaitTime(wait[step] * time.Second)
	url := "http://" + addr + "/updates"
	for {
		<-time.After(time.Duration(reportInterval) * time.Second)

		data, _ := json.Marshal(m.GetMetrics())

		var body bytes.Buffer
		gBody := gzip.NewWriter(&body)
		gBody.Write(data)
		gBody.Close()

		resp, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("Accept-Encoding", "gzip").
			SetHeader("Content-Encoding", "gzip").
			SetBody(&body).
			Post(url)

		if err != nil {
			if resp.StatusCode() == http.StatusGatewayTimeout {
				step++
				if step == 3 {
					logger.Log.Error("resty request error", zap.Error(err))
					return
				}
				client.SetRetryWaitTime(wait[step] * time.Second)
			}
			logger.Log.Error("resty request error", zap.Error(err))
			continue
		}

		step = 0

		logger.Log.Info("SendMetric", zap.String("URL", resp.Request.URL),
			zap.String("body", resp.String()),
			zap.String("StatusCode", resp.Status()))
	}
}
