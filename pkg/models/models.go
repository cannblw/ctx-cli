package models

import "time"

type State struct {
	ID       int64  `bun:"id,pk,autoincrement"`
	Name     string `bun:"name,notnull,unique"`
	Position int    `bun:"position,notnull"`
	Orphaned bool   `bun:"orphaned,default:false"`
}

type Context struct {
	ID          int64     `bun:"id,pk,autoincrement"`
	Name        string    `bun:"name,notnull,unique"`
	Description string    `bun:"description"`
	CreatedAt   time.Time `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt   time.Time `bun:"updated_at,notnull,default:current_timestamp"`
}

type Item struct {
	ID        int64     `bun:"id,pk,autoincrement"`
	Slug      string    `bun:"slug,notnull,unique"`
	ContextID *int64    `bun:"context_id"`
	Type      string    `bun:"type,notnull"`
	Value     string    `bun:"value,notnull"`
	StateID   *int64    `bun:"state_id"`
	Notes     string    `bun:"notes"`
	CreatedAt time.Time `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:"updated_at,notnull,default:current_timestamp"`

	Context *Context `bun:"rel:belongs-to,join:context_id=id"`
	State   *State   `bun:"rel:belongs-to,join:state_id=id"`
}
