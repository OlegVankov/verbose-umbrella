package monitor

import (
	"io"
	"log"
	"net/http"
	"time"
)

func SendMetrics(client *http.Client, m *Monitor, addr string, reportInterval int) {
	for {
		<-time.After(time.Duration(reportInterval) * time.Second)
		//for _, url := range m.GetRoutes(addr) {
		for _, body := range m.GetBody() {
			req, err := http.NewRequest(http.MethodPost, "http://"+addr+"/update/", body)
			req.Header.Set("Host", addr)
			req.Header.Set("Content-Type", "application/json")
			if err != nil {
				log.Println(err.Error())
				return
			}
			r, err := client.Do(req)
			if err != nil {
				log.Println(err.Error())
				return
			}
			rBody, _ := io.ReadAll(r.Body)
			log.Printf("URL: %s\n\tBody: %s\n\tStatusCode: %s\n", r.Request.URL, rBody, r.Status)
			_ = r.Body.Close()
		}
		// обнулим pollCounter после отправки метрик
		m.resetPollCount()
	}
}
