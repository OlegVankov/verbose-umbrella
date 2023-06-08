package storage

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMemStorage_UpdateGauge(t *testing.T) {
	type fields struct {
		Gauge   map[string]Gauge
		Counter map[string]Counter
	}
	type args struct {
		name string
		val  Gauge
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		wants  Gauge
	}{
		{name: "test #1", fields: fields{Gauge: map[string]Gauge{"alloc": 100}}, args: args{name: "alloc", val: 123}, wants: 123},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemStorage{
				Gauge:   tt.fields.Gauge,
				Counter: tt.fields.Counter,
			}
			m.UpdateGauge(tt.args.name, tt.args.val)
			assert.Equal(t, tt.wants, m.Gauge[tt.args.name], "значение не обновлено")
		})
	}
}

func TestMemStorage_UpdateCounter(t *testing.T) {
	type fields struct {
		Gauge   map[string]Gauge
		Counter map[string]Counter
	}
	type args struct {
		name string
		val  Counter
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		wants  Counter
	}{
		{name: "test #1", fields: fields{Counter: map[string]Counter{"count": 300}}, args: args{name: "count", val: 123}, wants: 423},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemStorage{
				Gauge:   tt.fields.Gauge,
				Counter: tt.fields.Counter,
			}
			m.UpdateCounter(tt.args.name, tt.args.val)
			assert.Equal(t, tt.wants, m.Counter[tt.args.name], "значение не увеличено")
		})
	}
}
