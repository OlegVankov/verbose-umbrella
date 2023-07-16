package pg

import (
	"context"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/OlegVankov/verbose-umbrella/internal/logger"
)

type Storage struct {
	db *sqlx.DB
}

func NewStorage(conn *sqlx.DB) *Storage {
	return &Storage{db: conn}
}

func (s *Storage) Bootstrap(ctx context.Context) error {
	query := `CREATE TABLE IF NOT EXISTS metrics 
(
	id varchar PRIMARY KEY,
	type varchar,
	delta bigint,
	value double precision
);
`
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = s.db.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Storage) UpdateGauge(ctx context.Context, id string, val float64) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Log.Error("gauge begin", zap.Error(err))
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `INSERT INTO metrics (id, type, value) VALUES ($1, $2, $3) ON CONFLICT (id) DO UPDATE SET value = $3`)
	if err != nil {
		logger.Log.Error("gauge prepare", zap.Error(err))
	}
	defer stmt.Close()

	stmt.QueryRowContext(ctx, id, "gauge", val)

	tx.Commit()
}

func (s *Storage) UpdateCounter(ctx context.Context, id string, val int64) (int64, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Log.Error("counter begin", zap.Error(err))
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `INSERT INTO metrics (id, type, delta) VALUES ($1, $2, $3) ON CONFLICT (id) DO UPDATE SET delta = $3`)
	if err != nil {
		logger.Log.Error("counter prepare", zap.Error(err))
	}
	defer stmt.Close()

	oldVal, _ := s.GetCounter(ctx, id)
	newVal := oldVal + val

	stmt.QueryRowContext(ctx, id, "counter", newVal)

	return newVal, tx.Commit()
}

func (s *Storage) GetGauge(ctx context.Context, id string) (float64, bool) {
	var val float64
	query := `SELECT value FROM metrics WHERE id = $1`
	err := s.db.QueryRowxContext(ctx, query, id).Scan(&val)
	if err != nil {
		return 0, false
	}
	return val, true
}

func (s *Storage) GetCounter(ctx context.Context, id string) (int64, bool) {
	var val int64
	query := `SELECT delta FROM metrics WHERE id = $1`
	err := s.db.QueryRowxContext(ctx, query, id).Scan(&val)
	if err != nil {
		return 0, false
	}
	return val, true
}

func (s *Storage) GetGaugeAll(ctx context.Context) map[string]float64 {
	// TODO implement me
	panic("implement me")
}

func (s *Storage) GetCounterAll(ctx context.Context) map[string]int64 {
	// TODO implement me
	panic("implement me")
}

func (s *Storage) PingStorage(ctx context.Context) error {
	return s.db.PingContext(ctx)
}
