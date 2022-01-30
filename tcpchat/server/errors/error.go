package errors

import "errors"

var (
	USER_NOT_EXIST     = errors.New("用户不存在")
	USER_ALREADY_EXIST = errors.New("用户已存在")
	INVALID_UID_PASSWD = errors.New("用户ID或密码错误")
)
