package gormutil

import (
	"context"

	"gorm.io/gorm"
)

type gormContextKey string

const txKey gormContextKey = "gorm_ctx_tx"

func InjectTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, txKey, tx)
}

/*
InjectTx 用于将事务注入到 context 中，以便在后续的数据库操作中使用
例如
func (uc *OrderUsecase) CreateOrder(ctx context.Context, order *Order, stockID int) error {
    return uc.db.Transaction(func(tx *gorm.DB) error {
        txCtx := injectTx(ctx, tx)

        // 订单保存 & 库存扣减，都走事务
        if err := uc.orderRepo.Save(txCtx, order); err != nil {
            return err
        }
        if err := uc.stockRepo.Deduct(txCtx, stockID, 1); err != nil {
            return err
        }
        return nil
    })
}
*/
