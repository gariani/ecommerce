package products

import (
	"context"

	repo "github.com/gariani/ecommerce/internal/adapters/postresql/sqlc"
)

type Service interface {
	ListProducts(ctx context.Context) ([]repo.Product, error)
	FindProductById(ctx context.Context, id int64) (repo.Product, error)
}

type svc struct {
	repo repo.Querier
}

func NewService(repo repo.Querier) Service {
	return &svc{repo: repo}
}

func (s *svc) ListProducts(ctx context.Context) ([]repo.Product, error) {
	return s.repo.ListProducts(ctx)
}

func (s *svc) FindProductById(ctx context.Context, id int64) (repo.Product, error) {
	return s.repo.FindProductbyID(ctx, id)
}
