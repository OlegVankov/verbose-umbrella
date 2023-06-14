package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemStorage_UpdateGauge(t *testing.T) {
	var ss Storage = NewStorage()
	ss.UpdateCounter("PollCount", 100)
	ss.UpdateGauge("alloc", 100)

	type args struct {
		name string
		val  Gauge
	}
	tests := []struct {
		name   string
		fields Storage
		args   args
		wants  Gauge
	}{
		{name: "test #1", fields: ss, args: args{name: "alloc", val: 123}, wants: 123},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.UpdateGauge(tt.args.name, tt.args.val)
			got, _ := tt.fields.GetGauge(tt.args.name)
			assert.Equal(t, tt.wants, got, "значение не обновлено")
		})
	}
}

func TestMemStorage_UpdateCounter(t *testing.T) {
	var ss Storage = NewStorage()
	ss.UpdateCounter("PollCount", 100)
	ss.UpdateGauge("alloc", 100)

	type args struct {
		name string
		val  Counter
	}
	tests := []struct {
		name   string
		fields Storage
		args   args
		wants  Counter
	}{
		{name: "test #1", fields: ss, args: args{name: "PollCount", val: 123}, wants: 223},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.UpdateCounter(tt.args.name, tt.args.val)
			got, _ := tt.fields.GetCounter(tt.args.name)
			assert.Equal(t, tt.wants, got, "значение не увеличено")
		})
	}
}

func TestMemStorage_GetGauge(t *testing.T) {
	var ss Storage = NewStorage()
	ss.UpdateCounter("PollCount", 100)
	ss.UpdateGauge("alloc", 100)

	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields Storage
		args   args
		want   Gauge
		want1  bool
	}{
		{name: "test #1", fields: ss, args: args{name: "alloc"}, want: 100, want1: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.fields.GetGauge(tt.args.name)
			assert.Equalf(t, tt.want, got, "GetGauge(%v)", tt.args.name)
			assert.Equalf(t, tt.want1, got1, "GetGauge(%v)", tt.args.name)
		})
	}
}

func TestMemStorage_GetCounter(t *testing.T) {
	var ss Storage = NewStorage()
	ss.UpdateCounter("PollCount", 100)
	ss.UpdateGauge("alloc", 100)
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields Storage
		args   args
		want   Counter
		want1  bool
	}{
		{name: "test #1", fields: ss, args: args{name: "PollCount"}, want: 100, want1: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, got1 := tt.fields.GetCounter(tt.args.name)
			assert.Equalf(t, tt.want, got, "GetCounter(%v)", tt.args.name)
			assert.Equalf(t, tt.want1, got1, "GetCounter(%v)", tt.args.name)
		})
	}
}
