package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func ToTimestamptzPtr(v *int64) pgtype.Timestamptz {
	if v == nil {
		return pgtype.Timestamptz{Valid: false}
	}
	return pgtype.Timestamptz{Time: time.UnixMilli(*v), Valid: true}
}

func ToUUIDPtr(v *uuid.UUID) pgtype.UUID {
	if v == nil {
		return pgtype.UUID{Valid: false}
	}
	return pgtype.UUID{Bytes: *v, Valid: true}
}

func ToInt64Ptr(t pgtype.Timestamptz) *int64 {
	if !t.Valid {
		return nil
	}
	v := t.Time.UnixMilli()
	return &v
}
