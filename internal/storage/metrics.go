package storage

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

//func (m *Metrics) UnmarshalJSON(bytes []byte) error {
//	fmt.Println("unmarshal :)")
//	json.Unmarshal(bytes, m)
//
//	return nil
//}
//
//func (m Metrics) MarshalJSON() ([]byte, error) {
//
//	fmt.Println("marshal :)")
//
//	return json.Marshal(m)
//}
