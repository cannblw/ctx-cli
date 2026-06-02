package models

type State struct {
	ID       int64  `bun:"id,pk,autoincrement"`
	Name     string `bun:"name,notnull,unique"`
	Position int    `bun:"position,notnull"`
	Orphaned bool   `bun:"orphaned,default:false"`
}
