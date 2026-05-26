// internal/repository/base.go
package repository

import (
	"context"

	"gorm.io/gorm"
)

type txKey struct{}

func GetDB(ctx context.Context, defaultDB *gorm.DB) *gorm.DB {
    if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
        return tx
    }
    return defaultDB
}

func RunInTransaction(db *gorm.DB, ctx context.Context, fn func(txCtx context.Context) error) error {
    return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        txCtx := context.WithValue(ctx, txKey{}, tx)
        return fn(txCtx)
    })
}