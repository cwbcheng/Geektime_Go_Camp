package dao

import (
	"database/sql"
	"github.com/pkg/errors" // 因为 dao 库不是第三方库，所以可以使用 errors
)

type UserInfo struct {}

func GetUserInfo(id int) (*UserInfo, error) {
	var result, err = queryDB(id)
	if err != nil {
		return nil, errors.Wrapf(err, "dao.GetUserInfo(%d) error", id)
	}

	return result, nil
}

func queryDB(id int) (*UserInfo, error) {
	// 连接数据库，构造查询语句查询，未查到数据返回 error
	return nil, sql.ErrNoRows
}