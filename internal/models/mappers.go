package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func IntPtrToDBTime(v *int64) pgtype.Timestamptz {
	if v == nil {
		return pgtype.Timestamptz{Valid: false}
	}
	return pgtype.Timestamptz{Time: time.UnixMilli(*v), Valid: true}
}

func UUIDPtrToDB(v *uuid.UUID) pgtype.UUID {
	if v == nil {
		return pgtype.UUID{Valid: false}
	}
	return pgtype.UUID{Bytes: *v, Valid: true}
}

func DBTimeToIntPtr(t pgtype.Timestamptz) *int64 {
	if !t.Valid {
		return nil
	}
	v := t.Time.UnixMilli()
	return &v
}

func DBTimeToInt(t pgtype.Timestamptz) int64 {
	if !t.Valid {
		return 0
	}
	return t.Time.UnixMilli()
}

func DBUUIDToUUIDPtr(v pgtype.UUID) *uuid.UUID {
	if !v.Valid {
		return nil
	}
	u := uuid.UUID(v.Bytes)
	return &u
}

func DBUUIDToUUID(v pgtype.UUID) uuid.UUID {
	if !v.Valid {
		return uuid.Nil
	}
	return v.Bytes
}
