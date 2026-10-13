package models

import (
	"strings"
	"testing"
	"time"
)

func TestSessionInputValidate(t *testing.T) {
	tz := time.FixedZone("MDT", -6*60*60)
	now := time.Date(2026, 10, 13, 12, 0, 0, 0, tz)
	library := &Location{ID: 1, GroupFriendly: true, Capacity: 8, OpenTime: "07:00", CloseTime: "23:00"}
	quietFloor := &Location{ID: 2, GroupFriendly: false, Capacity: 60, OpenTime: "07:00", CloseTime: "23:00"}

	valid := func() SessionInput {
		return SessionInput{
			LocationID:      1,
			StudentName:     "Ada Lovelace",
			Email:           "ada@byui.edu",
			StudentID:       "123456789",
			Date:            "2026-10-14",
			StartTime:       "15:00",
			DurationMinutes: 60,
			SessionType:     "individual",
		}
	}
	intPtr := func(n int) *int { return &n }

	tests := []struct {
		name      string
		mutate    func(*SessionInput)
		loc       *Location
		wantField string // empty means the input must be valid
	}{
		{name: "valid individual", mutate: func(*SessionInput) {}, loc: library},
		{name: "valid group", mutate: func(in *SessionInput) { in.SessionType = "group"; in.GroupSize = intPtr(4) }, loc: library},
		{name: "later today is fine", mutate: func(in *SessionInput) { in.Date = "2026-10-13"; in.StartTime = "18:00" }, loc: library},
		{name: "missing name", mutate: func(in *SessionInput) { in.StudentName = "" }, loc: library, wantField: "student_name"},
		{name: "bad email", mutate: func(in *SessionInput) { in.Email = "ada@" }, loc: library, wantField: "email"},
		{name: "email with display name", mutate: func(in *SessionInput) { in.Email = "Ada <ada@byui.edu>" }, loc: library, wantField: "email"},
		{name: "missing student id", mutate: func(in *SessionInput) { in.StudentID = "" }, loc: library, wantField: "student_id"},
		{name: "missing location", mutate: func(in *SessionInput) { in.LocationID = 0 }, loc: nil, wantField: "location_id"},
		{name: "unknown location", mutate: func(in *SessionInput) { in.LocationID = 99 }, loc: nil, wantField: "location_id"},
		{name: "past date", mutate: func(in *SessionInput) { in.Date = "2026-10-12" }, loc: library, wantField: "date"},
		{name: "earlier today", mutate: func(in *SessionInput) { in.Date = "2026-10-13"; in.StartTime = "09:00" }, loc: library, wantField: "start_time"},
		{name: "bad date format", mutate: func(in *SessionInput) { in.Date = "10/14/2026" }, loc: library, wantField: "date"},
		{name: "bad time format", mutate: func(in *SessionInput) { in.StartTime = "3pm" }, loc: library, wantField: "start_time"},
		{name: "duration not allowed", mutate: func(in *SessionInput) { in.DurationMinutes = 45 }, loc: library, wantField: "duration_minutes"},
		{name: "before opening", mutate: func(in *SessionInput) { in.StartTime = "06:30" }, loc: library, wantField: "start_time"},
		{name: "runs past closing", mutate: func(in *SessionInput) { in.StartTime = "22:30" }, loc: library, wantField: "start_time"},
		{name: "ends exactly at closing", mutate: func(in *SessionInput) { in.StartTime = "22:00" }, loc: library},
		{name: "unknown session type", mutate: func(in *SessionInput) { in.SessionType = "party" }, loc: library, wantField: "session_type"},
		{name: "group without size", mutate: func(in *SessionInput) { in.SessionType = "group" }, loc: library, wantField: "group_size"},
		{name: "group of one", mutate: func(in *SessionInput) { in.SessionType = "group"; in.GroupSize = intPtr(1) }, loc: library, wantField: "group_size"},
		{name: "group over capacity", mutate: func(in *SessionInput) { in.SessionType = "group"; in.GroupSize = intPtr(9) }, loc: library, wantField: "group_size"},
		{name: "group at quiet location", mutate: func(in *SessionInput) { in.SessionType = "group"; in.GroupSize = intPtr(3) }, loc: quietFloor, wantField: "session_type"},
		{name: "notes too long", mutate: func(in *SessionInput) { in.Notes = strings.Repeat("a", 501) }, loc: library, wantField: "notes"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := valid()
			tt.mutate(&in)
			in.Normalize()
			errs := in.Validate(tt.loc, now)

			if tt.wantField == "" {
				if len(errs) > 0 {
					t.Fatalf("expected valid input, got errors: %v", errs)
				}
				return
			}
			if _, ok := errs[tt.wantField]; !ok {
				t.Fatalf("expected an error on %q, got %v", tt.wantField, errs)
			}
		})
	}
}

func TestNormalizeDropsGroupSizeForIndividual(t *testing.T) {
	size := 4
	in := SessionInput{SessionType: " Individual ", GroupSize: &size, Email: " ADA@BYUI.EDU "}
	in.Normalize()
	if in.GroupSize != nil {
		t.Errorf("GroupSize = %v, want nil", *in.GroupSize)
	}
	if in.SessionType != "individual" || in.Email != "ada@byui.edu" {
		t.Errorf("got type %q email %q", in.SessionType, in.Email)
	}
}
