package postgres

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// pgxTimestamp converts time.Time to pgtype.Timestamptz for sqlc insert params.
func pgxTimestamp(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: t, Valid: true}
}
