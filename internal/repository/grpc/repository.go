package grpc

import "github.com/mephistolie/chefbook-backend-subscription/internal/config"

type Repository struct {
	Auth *Auth
}

func NewRepository(cfg *config.Config) (*Repository, error) {
	profileService, err := NewAuth(*cfg.AuthService.Addr)
	if err != nil {
		return nil, err
	}

	return &Repository{
		Auth: profileService,
	}, nil
}

func (r *Repository) Stop() error {
	_ = r.Auth.Conn.Close()
	return nil
}
