package pg

import (
	"context"
	"time"

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

func (s *Storage) UpdateGauge(ctx context.Context, id string, val float64) (err error) {
	for _, t := range []time.Duration{1, 3, 5} {
		if err = s.updateGauge(ctx, id, val); err == nil {
			return err
		}
		time.Sleep(t * time.Second)
	}
	return err
}

func (s *Storage) updateGauge(ctx context.Context, id string, val float64) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Log.Error("gauge begin", zap.Error(err))
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `INSERT INTO metrics (id, type, value) VALUES ($1, $2, $3) ON CONFLICT (id) DO UPDATE SET value = $3`)
	if err != nil {
		logger.Log.Error("gauge prepare", zap.Error(err))
		return err
	}
	defer stmt.Close()

	stmt.QueryRowContext(ctx, id, "gauge", val)

	return tx.Commit()
}

func (s *Storage) UpdateCounter(ctx context.Context, id string, val int64) (n int64, err error) {
	for _, t := range []time.Duration{1, 3, 5} {
		if n, err = s.updateCounter(ctx, id, val); err == nil {
			return n, err
		}
		time.Sleep(t * time.Second)
	}
	return 0, err
}

func (s *Storage) updateCounter(ctx context.Context, id string, val int64) (int64, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Log.Error("counter begin", zap.Error(err))
		return 0, err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `INSERT INTO metrics (id, type, delta) VALUES ($1, $2, $3) 
ON CONFLICT (id) DO UPDATE SET delta = metrics.delta + $3 RETURNING delta`)
	if err != nil {
		logger.Log.Error("counter prepare", zap.Error(err))
		return 0, err
	}
	defer stmt.Close()

	var newVal int64

	err = stmt.QueryRowContext(ctx, id, "counter", val).Scan(&newVal)
	if err != nil {
		logger.Log.Error("counter query", zap.Error(err))
		return 0, err
	}

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
	res := map[string]float64{}

	query := `SELECT id, value FROM metrics WHERE type = "gauge"`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		logger.Log.Error("get gauge all", zap.Error(err))
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		var value float64
		err = rows.Scan(&id, &value)
		if err != nil {
			return nil
		}
		res[id] = value
	}

	err = rows.Err()
	if err != nil {
		logger.Log.Error("get counter all", zap.Error(err))
		return nil
	}

	return res
}

func (s *Storage) GetCounterAll(ctx context.Context) map[string]int64 {
	res := map[string]int64{}

	query := `SELECT id, delta FROM metrics WHERE type = "counter"`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		logger.Log.Error("get counter all", zap.Error(err))
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		var delta int64
		err = rows.Scan(&id, &delta)
		if err != nil {
			return nil
		}
		res[id] = delta
	}

	err = rows.Err()
	if err != nil {
		logger.Log.Error("get counter all", zap.Error(err))
		return nil
	}

	return res
}

func (s *Storage) PingStorage(ctx context.Context) error {
	return s.db.PingContext(ctx)
}
