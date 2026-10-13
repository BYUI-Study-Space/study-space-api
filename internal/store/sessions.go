package store

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/jackc/pgx/v5"

	"github.com/BYUI-Study-Space/study-space-api/internal/models"
)

var uuidPattern = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

// CreateSession inserts an already validated session and returns it as saved.
func (s *Store) CreateSession(ctx context.Context, in models.SessionInput) (models.Session, error) {
	var id string
	err := s.db.QueryRow(ctx, `
		INSERT INTO sessions (location_id, student_name, email, student_id, session_date, start_time,
		                      duration_minutes, session_type, group_size, notes)
		VALUES ($1, $2, $3, $4, $5::date, $6::time, $7, $8, $9, NULLIF($10, ''))
		RETURNING id::text`,
		in.LocationID, in.StudentName, in.Email, in.StudentID, in.Date, in.StartTime,
		in.DurationMinutes, in.SessionType, in.GroupSize, in.Notes,
	).Scan(&id)
	if err != nil {
		return models.Session{}, fmt.Errorf("create session: %w", err)
	}
	return s.GetSession(ctx, id)
}

func (s *Store) GetSession(ctx context.Context, id string) (models.Session, error) {
	if !uuidPattern.MatchString(id) {
		return models.Session{}, ErrNotFound
	}
	// Column order must stay in sync with the models.Session fields.
	rows, err := s.db.Query(ctx, `
		SELECT s.id::text, s.location_id, l.name, s.student_name, s.email, s.student_id,
		       to_char(s.session_date, 'YYYY-MM-DD'),
		       to_char(s.start_time, 'HH24:MI'),
		       to_char(s.start_time + make_interval(mins => s.duration_minutes), 'HH24:MI'),
		       s.duration_minutes, s.session_type, s.group_size, COALESCE(s.notes, ''), s.created_at
		FROM sessions s
		JOIN locations l ON l.id = s.location_id
		WHERE s.id = $1`, id)
	if err != nil {
		return models.Session{}, fmt.Errorf("get session: %w", err)
	}
	session, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByPos[models.Session])
	if errors.Is(err, pgx.ErrNoRows) {
		return session, ErrNotFound
	}
	return session, err
}
