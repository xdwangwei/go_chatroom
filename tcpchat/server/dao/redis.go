package dao

import (
	"context"
	"github.com/go-redis/redis/v8"
	"runtime"
	"time"
)

var (
	ctx = context.Background()
	// 全局,对外提供
	redisClient *redis.Client
)

// 创建一个用于操作redis的客户端
func InitRedisClient() {
	redisClient = redis.NewClient(&redis.Options{
		// 网路类型, tcp 或 unix, 默认是 tcp
		Network: "tcp",
		// redis地址 , ip:port
		Addr: "127.0.0.1:6379",
		// 密码
		// Password: "",
		// 数据库
		DB: 0,
		// 可自定义连接函数,建立连接时调用,并且此参数配置优先于 前面两个参数
		//Dialer: func(ctx context.Context, network, addr string) (net.Conn, error) {
		//
		//},
		// 钩子函数,连接被建立时会被调用
		// 仅当客户端执行命令时需要从连接池获取连接时，若此时连接池需要新建连接时才会调用此钩子函数
		//OnConnect: func(ctx context.Context, conn *redis.Conn) error {
		//	fmt.Printf("一个连接建立, conn=%v\n", conn)
		//	return nil
		//},
		// 建立连接时最大重试次数,默认是3
		MaxRetries: 3,
		// 每次重试的最小间隔时间,默认是8ms,设置为-1则无间隔
		MinRetryBackoff: 8 * time.Millisecond,
		// 每次重试的最大间隔时间,默认是512ms,设置为-1则无间隔
		MaxRetryBackoff: 512 * time.Millisecond,

		// 连接池类型 true 表示 FIFO, false 表示 LIFO, FIFO 的开销更高
		PoolFIFO: false,
		// 最大连接数,默认是 runtime.GOMAXPROCS(给go设置的最大核心数) * 10
		// Default is 10 connections per every available CPU as reported by runtime.GOMAXPROCS.
		PoolSize: runtime.NumCPU() * 10,
		// 最小的闲置连接数
		MinIdleConns: 8,
		// 空闲连接的timeout,默认是5分钟
		IdleTimeout: 5 * time.Minute,
		// 空闲连接的检测时间间隔,默认是1分钟
		IdleCheckFrequency: time.Minute,

		// 拨号建立新连接时的timeout,默认是5s
		// Default is 5 seconds.
		DialTimeout: 5 * time.Second,
		// 数据包read的timeout,超过则命令执行失败
		// -1表示没有timeout,0表示默认,默认是3s
		// Default is 3 seconds.
		ReadTimeout: 3 * time.Second,
		// 数据包write的timeout,超过则命令执行失败
		// 默认和read的配置一致
		WriteTimeout: 3 * time.Second,

		// 如果连接全部被占用,客户请求的最大等待时间.默认是 前面设置的read的timeout+1
		// Default is ReadTimeout + 1 second.
		//PoolTimeout: time.Duration
		// Connection age at which client retires (closes) the connection.
		// Default is to not close aged connections.
		//MaxConnAge, time.Duration
	})

}
