package monitor

import (
	"log"
	"net/http"
	"time"
)

func SendMetrics(client *http.Client, m *Monitor, addr string, reportInterval int) {
	for {
		for _, url := range m.GetRoutes(addr) {
			req, err := http.NewRequest(http.MethodPost, url, nil)
			req.Header.Set("Host", addr)
			req.Header.Set("Content-Type", "text/plain")
			if err != nil {
				log.Println(err.Error())
				return
			}
			r, err := client.Do(req)
			if err != nil {
				log.Println(err.Error())
				return
			}
			log.Printf("StatusCode: %s\n", r.Status)
			r.Body.Close()
		}
		<-time.After(time.Duration(reportInterval) * time.Second)
	}
}
