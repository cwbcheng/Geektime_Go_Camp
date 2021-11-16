package service

import (
	pb "Geektime_Go_Camp/第四周/api"
	"Geektime_Go_Camp/第四周/internal/biz"
	"context"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
	biz *biz.UserBiz
}

func NewUserService(biz *biz.UserBiz) *UserService {
	return &UserService{
		biz: biz,
	}
}

func (s *UserService) UserInfo(ctx context.Context, in *pb.UserInfoRequest) (*pb.UserInfoReply, error) {
	u, err := s.biz.GetUserById(in.Uid)
	if err != nil {
		return nil, err
	} else {
		return &pb.UserInfoReply{User: &pb.User{
			Uid:   in.Uid,
			Name:  u.Nickname,
		}}, nil
	}

}