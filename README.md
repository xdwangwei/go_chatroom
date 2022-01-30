### 一、项目介绍

此项目（严格来说算不上项目）是在学习 golang 的过程中写的一个基于 tcp协议的简易聊天室，有兴趣可以看看，虽然比较小，但还是用到了很多东西。

### 二、关键技术

- gorouting （go 协程当然必不可少了，维持一个通信任务，接收啊发送消息什么的）
- go-redis （其实可以不用redis，但是为了结合所学知识，所以用redis来进行数据存储，主要是用于用于的登录/注册啊什么的）
- 借鉴了 MVC 模式，不过分离的不太好。
- 其他好像没什么值得说的了吧

### 三、实现功能

- 客户注册、登录（不允许重复登录，重复注册，都有校验以及错误提示）
- 在线用户列表查询（可以看到当前系统有哪些人在线，除了自己）
- 群聊消息（发送一条消息，当然在线的人都能收到，除了自己）
- 私聊消息（可以私聊某一个人，若这个人不存在或者不在线，也会有相应提示）
- 退出系统/退出登录（退出登录后，其他人的在线列表会被更新，此时再次查看在线列表就看不到我啦）
- 待续....

### 四、注意事项

- 下载源代码后，请修改 redis.go 中redis的连接信息。也就是 host:port， 记得改为你自己的。

- 记得先安装 [go-redis](https://github.com/go-redis/redis)，我这里是通过 go mod 装的，就是 进入 GOPATH 目录下，执行下面这个命令

  ```sh
  go get github.com/go-redis/redis/v8
  ```

- 记得把文件夹放入你的 GOPATH/src 路径下，项目可以通过 goland 打开，goland 能够直接运行，所以挺方便的。

- 更好的方式是分别编译客户端和服务端代码，然后得到可执行文件，就可以打开cmd窗口运行啦，此时可以启动多个客户端哦。

  ```sh
  # 进入 GOPATH/src/awesomeProject，分别执行
  go build -o server.exe .\tcpchat\server\main\
  go build -o client.exe .\tcpchat\client\main\
  ```

### 五、运行截图

![image-20220130202434734](http://typora.iwangwei.top/img/image-20220130202434734.png)