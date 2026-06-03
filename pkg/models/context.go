package models

import "time"

type Context struct {
	ID          int64     `bun:"id,pk,autoincrement"`
	Name        string    `bun:"name,notnull,unique"`
	Description string    `bun:"description"`
	CreatedAt   time.Time `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt   time.Time `bun:"updated_at,notnull,default:current_timestamp"`
}
