package repository

import (
	"backend/src/internal/repository/postgre"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewLinksRepository(storageType string, conn *pgxpool.Pool) (ILinksRepository, error) {
	var repo ILinksRepository
	switch storageType {
	case "postgres":
		repo = postgre.NewLinksRepository(conn)
	case "local":
	default:
		return nil, errors.New("invalid storage type")
	}
	return repo, nil
}
