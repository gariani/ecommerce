package orders

import (
	"context"
	"errors"
	"fmt"

	repo "github.com/gariani/ecommerce/internal/adapters/postresql/sqlc"
	"github.com/jackc/pgx/v5"
)

var (
	ErrorProductNotFound = errors.New(`product not found`)
	ErrorProductNoStock  = errors.New(`product not in stock`)
)

type svc struct {
	repo *repo.Queries
	db   *pgx.Conn
}

func NewService(repo *repo.Queries, db *pgx.Conn) Service {
	return &svc{
		repo: repo,
		db:   db,
	}
}

func (s *svc) PlaceOrder(ctx context.Context, tempOrder createOrderParams) (repo.Order, error) {

	if tempOrder.CustomerID == 0 {
		return repo.Order{}, fmt.Errorf("customer ID is required")
	}

	if len(tempOrder.Items) == 0 {
		return repo.Order{}, fmt.Errorf("")
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return repo.Order{}, err
	}

	defer tx.Rollback(ctx)

	qtx := s.repo.WithTx(tx)

	order, err := qtx.CreateOrder(ctx, tempOrder.CustomerID)

	if err != nil {
		return repo.Order{}, err
	}

	for _, item := range tempOrder.Items {
		product, err := s.repo.FindProductbyID(ctx, item.ProductID)
		if err != nil {
			return repo.Order{}, ErrorProductNotFound
		}

		if product.Quantity < item.Quantity {
			return repo.Order{}, ErrorProductNoStock
		}

		_, err = qtx.CreateOrderItem(ctx, repo.CreateOrderItemParams{
			OrderID:    order.ID,
			ProductID:  item.ProductID,
			Quantity:   item.Quantity,
			PriceCents: product.PriceInCenters,
		})

		if err != nil {
			return repo.Order{}, err
		}
	}

	tx.Commit(ctx)

	return order, nil

}
