package dao

import (
	"awesomeProject/tcpchat/server/errors"
	"awesomeProject/tcpchat/server/model"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
)

// 构造全局变量
var (
	GlobalUserDAO *UserDAO
)

type UserDAO struct {
}

// 初始化
func init() {
	if GlobalUserDAO == nil {
		GlobalUserDAO = &UserDAO{}
	}
}

// 判断用户是否存在
func (this *UserDAO) GetUserById(userId int) (str string, err error) {
	str, err = redisClient.HGet(ctx, "users", strconv.Itoa(userId)).Result()
	// 用户不存在
	if err == redis.Nil {
		err = errors.USER_NOT_EXIST
		return
	}
	return
}

// 登录--从redis进行校验
func (this *UserDAO) Login(userId int, password string) (user model.User, err error) {

	str, err := redisClient.HGet(ctx, "users", strconv.Itoa(userId)).Result()
	// 用户不存在
	if err == redis.Nil {
		err = errors.USER_NOT_EXIST
		return
	}
	err = json.Unmarshal([]byte(str), &user)
	if err != nil {
		fmt.Errorf("json unmarshal error[%v]!\n", err)
		return
	}
	// 密码比对失败
	if password != user.Password {
		err = errors.INVALID_UID_PASSWD
		return
	}
	// 登录成功
	return
}

// 注册到redis
func (this *UserDAO) Register(user *model.User) (err error) {
	str, err := redisClient.HGet(ctx, "users", strconv.Itoa(user.UserId)).Result()
	// 用户ID已存在
	if err != redis.Nil || str != "" {
		err = errors.USER_ALREADY_EXIST
		return
	}
	// 注册
	bytes, err := json.Marshal(user)
	if err != nil {
		fmt.Errorf("json marshal error[%v]!\n", err)
		return
	}
	err = redisClient.HSet(ctx, "users", strconv.Itoa(user.UserId), string(bytes)).Err()
	// 用户不存在
	if err != nil {
		fmt.Errorf("redis command [hset] error[%v]!\n", err)
		return
	}
	// 注册成功
	return
}
