package bootsrap

import (
	"unisched/server/db"
	"unisched/server/internal/repository"
	"unisched/server/internal/services"
)

func StartServices(dbEngine *db.DBEngine) (*services.Services, error) {
	repo, err := repository.NewRepositories(dbEngine.Client)
	if err != nil {
		return nil, err
	}
	return services.NewServices(repo), nil
}
