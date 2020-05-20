package routers

import (
	"github.com/1024casts/snake/handler"

	"github.com/gin-gonic/gin"

	"github.com/swaggo/gin-swagger" //nolint: goimports
	"github.com/swaggo/gin-swagger/swaggerFiles"

	// import swagger handler
	_ "github.com/1024casts/snake/docs" // docs is generated by Swag CLI, you have to import it.
	"github.com/1024casts/snake/handler/user"
	"github.com/1024casts/snake/router/middleware"
)

// Load loads the middlewares, routes, handlers.
func Load(g *gin.Engine, mw ...gin.HandlerFunc) *gin.Engine {
	// 使用中间件
	g.Use(middleware.NoCache)
	g.Use(middleware.Options)
	g.Use(middleware.Secure)
	g.Use(middleware.Logging())
	g.Use(middleware.RequestID())
	g.Use(mw...)

	// 404 Handler.
	g.NoRoute(handler.RouteNotFound)
	g.NoMethod(handler.RouteNotFound)

	// 静态资源，主要是图片
	g.Static("/static", "./static")

	// swagger api docs
	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// pprof router 性能分析路由
	// 默认关闭，开发环境下可以打开
	// 访问方式: HOST/debug/pprof
	// 通过 HOST/debug/pprof/profile 生成profile
	// 查看分析图 go tool pprof -http=:5000 profile
	// pprof.Register(g)

	// 下面就可以开始写具体的业务路由了

	// api for authentication functionalities
	g.POST("/v1/register", user.Register)
	g.POST("/v1/login", user.Login)
	g.GET("/v1/vcode", user.VCode)
	// 手机号登录
	g.POST("/v1/login/phone", user.PhoneLogin)

	// The user handlers, requiring authentication
	g.GET("/v1/users/:id", user.Get)
	u := g.Group("/v1/users")
	u.Use(middleware.AuthMiddleware())
	{
		u.PUT("/:id", user.Update)
	}

	return g
}
