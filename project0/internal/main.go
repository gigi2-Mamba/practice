package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"project0/internal/repository"
	"project0/internal/repository/dao"
	"project0/internal/repository/mysqlInit"
	"project0/internal/service"
	"project0/internal/web"
	"project0/internal/web/middlewares"
	"strings"
	"time"
)

func initUserHdl(db *gorm.DB, server *gin.Engine) {
	ud := dao.NewUserDAO(db)
	ur := repository.NewUserRepository(ud)
	us := service.NewUserService(ur)
	user := web.NewUserHandler(us)
	user.RegisterRoute(server)
}

func main() {
	// gorm连接mysql
	err := mysqlInit.Init()
	if err != nil {
		return
	}
	log.Println("mysqlInit init success")
	server := gin.Default()

	server.Use(cors.New(cors.Config{
		AllowCredentials: true, // 允许cookie
		AllowHeaders:     []string{"Content-Type", "authorization"},
		//AllowOrigins:     []string{"http://localhost:3000"},
		MaxAge: 12 * time.Hour,
		AllowOriginFunc: func(origin string) bool {
			log.Println(strings.HasPrefix(origin, "http://localhost"))
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "company.com")
		},
	}), func(context *gin.Context) {
		log.Println("work? ")
	})
	login := &middlewares.LoginMiddlewareBuilder{}
	// 先初始化 存储数据,直接存cookie作教学
	store := cookie.NewStore([]byte("secret"))

	server.Use(sessions.Sessions("ssid", store), login.CheckLogin())

	// 注册跨域中间件
	//server.Use()

	//ud := dao.NewUserDAO(mysqlInit.Db)
	//ur := repository.NewUserRepository(ud)
	//us := service.NewUserService(ur)
	//user := web.NewUserHandler(us)
	//user.RegisterRoute(server)
	initUserHdl(mysqlInit.Db, server)
	server.Run(":8083")
}
