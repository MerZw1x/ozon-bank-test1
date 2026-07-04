package repository

func NewLinksRepository(storageType string, dbdsn string) (*ILinksRepository, error) {
	switch storageType {
	case "postgres":
	case "local"
	}
}
