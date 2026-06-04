package entities

import "time"

type ProductImage struct {
	ID        int64
	ProductID int64
	URL       string
	IsMain    bool
	CreatedAt time.Time
}
