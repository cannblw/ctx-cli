package repo

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"

	_ "github.com/cannblw/ctx-cli/migrations"

	"github.com/pressly/goose/v3"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"

	"github.com/cannblw/ctx-cli/pkg/models"
)

type BunRepository struct {
	db *bun.DB
}

func NewBunRepository(dbPath string) (*BunRepository, error) {
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("create db dir: %w", err)
	}

	sqldb, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	goose.SetDialect("sqlite3")
	if err := goose.Up(sqldb, findMigrationsDir()); err != nil {
		sqldb.Close()
		return nil, fmt.Errorf("run migrations: %w", err)
	}

	db := bun.NewDB(sqldb, sqlitedialect.New())
	return &BunRepository{db: db}, nil
}

func NewBunRepositoryFromDB(db *bun.DB) *BunRepository {
	return &BunRepository{db: db}
}

func (r *BunRepository) CreateContext(ctx context.Context, name, description string) (*models.Context, error) {
	c := &models.Context{
		Name:        name,
		Description: description,
	}
	_, err := r.db.NewInsert().Model(c).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("create context: %w", err)
	}
	return c, nil
}

func (r *BunRepository) GetContext(ctx context.Context, name string) (*models.Context, error) {
	c := new(models.Context)
	err := r.db.NewSelect().Model(c).Where("name = ?", name).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("get context %q: %w", name, err)
	}
	return c, nil
}

func (r *BunRepository) ListContexts(ctx context.Context) ([]models.Context, error) {
	var contexts []models.Context
	err := r.db.NewSelect().Model(&contexts).OrderExpr("created_at ASC").Scan(ctx)
	return contexts, err
}

func (r *BunRepository) Close() error {
	return r.db.Close()
}

func findMigrationsDir() string {
	dirs := []string{"migrations", "../migrations", "../../migrations", "../../../migrations"}
	for _, d := range dirs {
		if fi, err := os.Stat(d); err == nil && fi.IsDir() {
			abspath, _ := filepath.Abs(d)
			return abspath
		}
	}
	return "migrations"
}
