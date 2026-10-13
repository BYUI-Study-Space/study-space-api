// Package store runs the SQL queries against Postgres.
package store

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/BYUI-Study-Space/study-space-api/internal/models"
)

var ErrNotFound = errors.New("not found")

type Store struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Store {
	return &Store{db: db}
}

func (s *Store) Ping(ctx context.Context) error {
	return s.db.Ping(ctx)
}

// locationColumns must stay in the same order as the models.Location fields.
const locationColumns = `id, name, building, room, description, image_url, has_wifi, group_friendly,
	noise_level, capacity, to_char(open_time, 'HH24:MI'), to_char(close_time, 'HH24:MI')`

func (s *Store) ListLocations(ctx context.Context, f models.LocationFilter) ([]models.Location, error) {
	var where []string
	var args []any
	add := func(cond string, arg any) {
		args = append(args, arg)
		where = append(where, fmt.Sprintf(cond, len(args)))
	}
	if f.Wifi != nil {
		add("has_wifi = $%d", *f.Wifi)
	}
	if f.Group != nil {
		add("group_friendly = $%d", *f.Group)
	}
	if f.Noise != "" {
		add("noise_level = $%d", f.Noise)
	}

	query := "SELECT " + locationColumns + " FROM locations"
	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
	}
	query += " ORDER BY name"

	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list locations: %w", err)
	}
	return pgx.CollectRows(rows, pgx.RowToStructByPos[models.Location])
}

func (s *Store) GetLocation(ctx context.Context, id int64) (models.Location, error) {
	rows, err := s.db.Query(ctx, "SELECT "+locationColumns+" FROM locations WHERE id = $1", id)
	if err != nil {
		return models.Location{}, fmt.Errorf("get location: %w", err)
	}
	loc, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByPos[models.Location])
	if errors.Is(err, pgx.ErrNoRows) {
		return loc, ErrNotFound
	}
	return loc, err
}
