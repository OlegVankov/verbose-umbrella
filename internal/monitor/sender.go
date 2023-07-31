package monitor

import (
	"bytes"
	"compress/gzip"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"sync"
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

func SendBatch(m *Monitor, addr string, reportInterval int, key string, wg *sync.WaitGroup) {
	defer wg.Done()
	wait := []time.Duration{1, 3, 5}
	step := 0

	client := resty.New()
	// client.SetRetryCount(3)
	client.SetRetryWaitTime(wait[step] * time.Second)

	url := "http://" + addr + "/updates"
	for {
		<-time.After(time.Duration(reportInterval) * time.Second)

		data, err := json.Marshal(m.GetMetrics())
		if err != nil {
			logger.Log.Error("send batch", zap.Error(err))
			continue
		}

		var body bytes.Buffer
		gBody := gzip.NewWriter(&body)
		gBody.Write(data)
		gBody.Close()

		sha := computeHmac256(data, key)

		resp, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("Accept-Encoding", "gzip").
			SetHeader("Content-Encoding", "gzip").
			SetHeader("HashSHA256", sha).
			SetBody(&body).
			Post(url)

		if err != nil {
			if resp.StatusCode() == http.StatusGatewayTimeout {
				step++
				if step == 3 {
					logger.Log.Error("send batch", zap.Error(err))
					return
				}
				client = client.SetRetryWaitTime(wait[step] * time.Second)
			}
			logger.Log.Error("send batch", zap.Error(err))
			continue
		}

		step = 0

		logger.Log.Info("SendMetric", zap.String("URL", resp.Request.URL),
			zap.String("body", resp.String()),
			zap.String("StatusCode", resp.Status()))
	}
}

func computeHmac256(message []byte, secret string) string {
	if secret == "" {
		return ""
	}
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write(message)
	return hex.EncodeToString(h.Sum(nil))
}
