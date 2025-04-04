package user

import (
	"context"
	"fmt"
)

func NewHttpHandler(ctx context.Context) *HttpHandler {
	return &HttpHandler{service: New(ctx)}
}

type Service struct {
	ctx context.Context
}

func New(ctx context.Context) *Service {
	return &Service{ctx: ctx}
}

func (s *Service) GetUser(id string) string {
	return fmt.Sprintf("User(ID=%s)", id)
}

func (s *Service) CreateUser(name string) string {
	return fmt.Sprintf("User(Name=%s)", name)
}
