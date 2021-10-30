package dao

import (
	"database/sql"
	"github.com/pkg/errors" // 因为 dao 库不是第三方库，所以可以使用 errors
)

type UserInfo struct {}

func GetUserInfo(id int) (*UserInfo, error) {
	var result, err = queryDB(id)
	if err != nil {
		// 应该把 error 包装抛给上层，让业务逻辑判断根因，然后做具体处理。
		// 比如本该有的信息没有查到，就返回 error；
		// 比如没有用户信息就调用其他方法创建一个新的。这样就是把这个 error 处理掉不继续抛出 error。
		return nil, errors.Wrapf(err, "dao.GetUserInfo(%d) error", id)
	}

	return result, nil
}

func queryDB(id int) (*UserInfo, error) {
	// 连接数据库，构造查询语句查询，未查到数据返回 error
	return nil, sql.ErrNoRows
}