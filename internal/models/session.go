package models

import (
	"fmt"
	"net/mail"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	dateLayout = "2006-01-02"
	timeLayout = "15:04"

	maxNameLength      = 100
	maxStudentIDLength = 20
	maxNotesLength     = 500
)

var allowedDurations = map[int]bool{30: true, 60: true, 90: true, 120: true, 180: true}

// SessionInput is the body of POST /api/sessions.
type SessionInput struct {
	LocationID      int64  `json:"location_id"`
	StudentName     string `json:"student_name"`
	Email           string `json:"email"`
	StudentID       string `json:"student_id"`
	Date            string `json:"date"`       // "YYYY-MM-DD"
	StartTime       string `json:"start_time"` // "HH:MM"
	DurationMinutes int    `json:"duration_minutes"`
	SessionType     string `json:"session_type"` // "individual" or "group"
	GroupSize       *int   `json:"group_size"`
	Notes           string `json:"notes"`
}

// Session is a saved session. Field order matches the store's SELECT column order.
type Session struct {
	ID              string    `json:"id"`
	LocationID      int64     `json:"location_id"`
	LocationName    string    `json:"location_name"`
	StudentName     string    `json:"student_name"`
	Email           string    `json:"email"`
	StudentID       string    `json:"student_id"`
	Date            string    `json:"date"`
	StartTime       string    `json:"start_time"`
	EndTime         string    `json:"end_time"`
	DurationMinutes int       `json:"duration_minutes"`
	SessionType     string    `json:"session_type"`
	GroupSize       *int      `json:"group_size"`
	Notes           string    `json:"notes"`
	CreatedAt       time.Time `json:"created_at"`
}

// Normalize trims whitespace and drops group_size for individual sessions.
func (in *SessionInput) Normalize() {
	in.StudentName = strings.TrimSpace(in.StudentName)
	in.Email = strings.ToLower(strings.TrimSpace(in.Email))
	in.StudentID = strings.TrimSpace(in.StudentID)
	in.Date = strings.TrimSpace(in.Date)
	in.StartTime = strings.TrimSpace(in.StartTime)
	in.SessionType = strings.ToLower(strings.TrimSpace(in.SessionType))
	in.Notes = strings.TrimSpace(in.Notes)
	if in.SessionType == "individual" {
		in.GroupSize = nil
	}
}

// Validate returns a message per invalid field, keyed by its JSON name; an empty map means valid.
// loc is the requested location, or nil when it doesn't exist. now must be in the campus timezone.
func (in *SessionInput) Validate(loc *Location, now time.Time) map[string]string {
	errs := map[string]string{}

	switch {
	case in.StudentName == "":
		errs["student_name"] = "is required"
	case utf8.RuneCountInString(in.StudentName) > maxNameLength:
		errs["student_name"] = fmt.Sprintf("must be at most %d characters", maxNameLength)
	}

	if in.Email == "" {
		errs["email"] = "is required"
	} else if addr, err := mail.ParseAddress(in.Email); err != nil || addr.Address != in.Email {
		errs["email"] = "must be a valid email"
	}

	switch {
	case in.StudentID == "":
		errs["student_id"] = "is required"
	case len(in.StudentID) > maxStudentIDLength:
		errs["student_id"] = fmt.Sprintf("must be at most %d characters", maxStudentIDLength)
	}

	switch {
	case in.LocationID == 0:
		errs["location_id"] = "is required"
	case loc == nil:
		errs["location_id"] = "unknown location"
	}

	if !allowedDurations[in.DurationMinutes] {
		errs["duration_minutes"] = "must be 30, 60, 90, 120 or 180"
	}

	switch in.SessionType {
	case "individual":
	case "group":
		switch {
		case loc != nil && !loc.GroupFriendly:
			errs["session_type"] = "this location does not allow group sessions"
		case in.GroupSize == nil:
			errs["group_size"] = "is required for group sessions"
		case *in.GroupSize < 2:
			errs["group_size"] = "must be at least 2"
		case loc != nil && *in.GroupSize > loc.Capacity:
			errs["group_size"] = fmt.Sprintf("exceeds this location's capacity of %d", loc.Capacity)
		}
	case "":
		errs["session_type"] = "is required"
	default:
		errs["session_type"] = `must be "individual" or "group"`
	}

	if utf8.RuneCountInString(in.Notes) > maxNotesLength {
		errs["notes"] = fmt.Sprintf("must be at most %d characters", maxNotesLength)
	}

	in.validateSchedule(loc, now, errs)
	return errs
}

func (in *SessionInput) validateSchedule(loc *Location, now time.Time, errs map[string]string) {
	date, dateErr := time.ParseInLocation(dateLayout, in.Date, now.Location())
	switch {
	case in.Date == "":
		errs["date"] = "is required"
	case dateErr != nil:
		errs["date"] = "must be a date in YYYY-MM-DD format"
	}

	start, startErr := time.Parse(timeLayout, in.StartTime)
	switch {
	case in.StartTime == "":
		errs["start_time"] = "is required"
	case startErr != nil:
		errs["start_time"] = "must be a time in HH:MM format"
	}

	if dateErr != nil || startErr != nil {
		return
	}

	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	startsAt := time.Date(date.Year(), date.Month(), date.Day(), start.Hour(), start.Minute(), 0, 0, now.Location())
	switch {
	case date.Before(today):
		errs["date"] = "cannot be in the past"
	case startsAt.Before(now):
		errs["start_time"] = "cannot be in the past"
	}

	if loc == nil || errs["duration_minutes"] != "" {
		return
	}
	open, err1 := time.Parse(timeLayout, loc.OpenTime)
	closing, err2 := time.Parse(timeLayout, loc.CloseTime)
	if err1 != nil || err2 != nil {
		return
	}
	startMin := minuteOfDay(start)
	if startMin < minuteOfDay(open) || startMin+in.DurationMinutes > minuteOfDay(closing) {
		errs["start_time"] = fmt.Sprintf("session must fit within opening hours (%s–%s)", loc.OpenTime, loc.CloseTime)
	}
}

func minuteOfDay(t time.Time) int {
	return t.Hour()*60 + t.Minute()
}
