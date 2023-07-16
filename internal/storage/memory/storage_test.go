package memory

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/OlegVankov/verbose-umbrella/internal/storage"
)

var ctx = context.Background()

func TestMemStorage_UpdateGauge(t *testing.T) {
	ss := NewStorage()
	ss.UpdateCounter(ctx, "PollCount", 100)
	ss.UpdateGauge(ctx, "alloc", 100)

	type args struct {
		name string
		val  float64
	}
	tests := []struct {
		name   string
		fields storage.Storage
		args   args
		wants  float64
	}{
		{name: "test #1", fields: ss, args: args{name: "alloc", val: 123}, wants: 123},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.UpdateGauge(ctx, tt.args.name, tt.args.val)
			got, _ := tt.fields.GetGauge(ctx, tt.args.name)
			assert.Equal(t, tt.wants, got, "значение не обновлено")
		})
	}
}

func TestMemStorage_UpdateCounter(t *testing.T) {
	ss := NewStorage()
	ss.UpdateCounter(ctx, "PollCount", 100)
	ss.UpdateGauge(ctx, "alloc", 100)

	type args struct {
		name string
		val  int64
	}
	tests := []struct {
		name   string
		fields storage.Storage
		args   args
		wants  int64
	}{
		{name: "test #1", fields: ss, args: args{name: "PollCount", val: 123}, wants: 223},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.UpdateCounter(ctx, tt.args.name, tt.args.val)
			got, _ := tt.fields.GetCounter(ctx, tt.args.name)
			assert.Equal(t, tt.wants, got, "значение не увеличено")
		})
	}
}

func TestMemStorage_GetGauge(t *testing.T) {
	ss := NewStorage()
	ss.UpdateCounter(ctx, "PollCount", 100)
	ss.UpdateGauge(ctx, "alloc", 100)

	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields storage.Storage
		args   args
		want   float64
		want1  bool
	}{
		{name: "test #1", fields: ss, args: args{name: "alloc"}, want: 100, want1: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.fields.GetGauge(ctx, tt.args.name)
			assert.Equalf(t, tt.want, got, "GetGauge(%v)", tt.args.name)
			assert.Equalf(t, tt.want1, got1, "GetGauge(%v)", tt.args.name)
		})
	}
}

func TestMemStorage_GetCounter(t *testing.T) {
	var ss storage.Storage = NewStorage()
	ss.UpdateCounter(ctx, "PollCount", 100)
	ss.UpdateGauge(context.Background(), "alloc", 100)
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields storage.Storage
		args   args
		want   int64
		want1  bool
	}{
		{name: "test #1", fields: ss, args: args{name: "PollCount"}, want: 100, want1: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, got1 := tt.fields.GetCounter(ctx, tt.args.name)
			assert.Equalf(t, tt.want, got, "GetCounter(%v)", tt.args.name)
			assert.Equalf(t, tt.want1, got1, "GetCounter(%v)", tt.args.name)
		})
	}
}
