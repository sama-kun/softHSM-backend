package models

import "time"

type BaseModel struct {
	Id        int64      `json:"id"`
	IsDeleted bool       `json:"isDeleted,omitempty"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
	UpdatedAt time.Time  `json:"updatedAt,omitempty"`
	CreatedAt time.Time  `json:"createdAt,omitempty"`
}
