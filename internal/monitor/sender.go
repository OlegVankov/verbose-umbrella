package monitor

import (
	"time"

	"github.com/OlegVankov/verbose-umbrella/internal/logger"
	"github.com/OlegVankov/verbose-umbrella/internal/storage"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
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

			//logger.Log.Info("Header", zap.String("content-encoding", resp.Header().Get("Content-Encoding")))

			//if resp.Header().Get("Content-Encoding") == "gzip" {
			//	g, err := gzip.NewReader(resp.RawBody())
			//	if err != nil {
			//		logger.Log.Error("gzip reader error", zap.Error(err))
			//		continue
			//	}
			//	json.NewDecoder(g).Decode(&metric)
			//	g.Close()
			//}

			logger.Log.Info("SendMetric", zap.String("URL", resp.Request.URL),
				zap.String("body", resp.String()),
				zap.String("StatusCode", resp.Status()),
				zap.Any("metric", metric))
		}
		// обнулим pollCounter после отправки метрик
		m.resetPollCount()
	}
}
