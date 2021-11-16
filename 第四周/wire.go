//+build wireinject

package main

import (
	"Geektime_Go_Camp/第四周/internal/biz"
	"Geektime_Go_Camp/第四周/internal/data"
	"github.com/google/wire"
	"Geektime_Go_Camp/第四周/internal/service"
)

func InitUserService() *service.UserService {
	wire.Build(service.NewUserService, biz.NewUserBiz, data.NewUserRepo, data.NewDB)
	return &service.UserService{}
}