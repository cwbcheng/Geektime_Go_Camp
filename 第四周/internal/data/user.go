package data

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (u *UserRepo) GetUser(id uint64) (*UserPO, error) {
	//查询数据库
	return &UserPO{
		Nickname: "Tom",
	}, nil
}

type UserPO struct {
	Nickname string
}


type UserDO struct {
	Nickname string
}

func convert()  {
	bytes, _ := json.Marshal(&UserPO{})
	do := &UserDO{}
	json.Unmarshal(bytes, do)
}