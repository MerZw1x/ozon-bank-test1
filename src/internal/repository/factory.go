package repository

import (
	"backend/src/internal/repository/local"
	"backend/src/internal/repository/postgre"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewLinksRepository(storageType string, conn *pgxpool.Pool) (ILinksRepository, error) {
	switch storageType {
	case "postgres":
		return postgre.NewLinksRepository(conn), nil
	case "local":
		return local.NewLinksRepository(), nil
	default:
		return nil, errors.New("invalid storage type")
	}
}
