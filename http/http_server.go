package http

import (
	"app/build"
	"app/config"
	"app/docs"
	"app/dto"
	"app/errors"
	"app/pkg/apperror"
	"app/pkg/logger"
	"app/pkg/utils"
	"app/service"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

type Server struct {
	AccountSvc    *service.AccountService
	AuthSvc       *service.AuthService
	UserSvc       *service.UserService
	ChatSvc       *service.ChatService
	GroupSvc      *service.GroupService
	SocketHandler *service.WebSocketHandler
}

func NewServer(accountSvc *service.AccountService, authSvc *service.AuthService, userSvc *service.UserService, chatSvc *service.ChatService, groupSvc *service.GroupService, socketHandler *service.WebSocketHandler) *Server {
	return &Server{AccountSvc: accountSvc, AuthSvc: authSvc, UserSvc: userSvc, ChatSvc: chatSvc, GroupSvc: groupSvc, SocketHandler: socketHandler}
}

func (s *Server) Routes(router *gin.RouterGroup) {
	router.GET("/health", func(ctx *gin.Context) {
		ctx.String(200, build.Info().String())
	})
	if !config.Cfg.HTTP.IsProduction {
		router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	router.POST("/auth/register", s.Register)
	router.POST("/auth/login", s.Login)
	router.POST("/auth/logout", s.Logout)

	router.GET("/account/profile", s.Authenticate, s.GetProfile)
	router.GET("/account/profile/:phoneNumber", s.Authenticate, s.GetProfileByPhoneNumber)
	router.GET("/account/profile/suggest", s.Authenticate, s.GetSuggestFriendProfiles)
	router.GET("/account/profile/userID/:userID", s.Authenticate, s.GetProfileByUserID)
	router.GET("/account/info", s.Authenticate, s.GetAccountProfile)
	router.POST("/account/check-phone", s.CheckPhoneNumber)
	router.PUT("/account/reset-password", s.ResetPassword)
	router.PUT("/account/change-password", s.Authenticate, s.ChangePassword)
	router.PUT("/account/change-avatar", s.Authenticate, s.ChangeAvatar)
	router.PUT("/account/change-profile", s.Authenticate, s.ChangeProfile)

	router.POST("/user/create", s.CreateUser)
	router.GET("/user/info/:id", s.Authenticate, s.GetUser)
	router.GET("/user/update-avatar-account", s.Authenticate, s.UpdateAvatarAsync)

	router.POST("/chat/create", s.Authenticate, s.CreateChat)
	router.GET("/chat/x-to-y", s.Authenticate, s.GetChatActivityFromNToM)
	router.GET("/chat/search-bkw", s.Authenticate, s.SearchByKeyWord)
	router.GET("/chat/get-search", s.Authenticate, s.SearchByKeyWord)

	router.POST("/group/create", s.Authenticate, s.CreateGroup)
	router.GET("/group/info", s.Authenticate, s.GetGroupInfo)
}

func (s *Server) Socket(router *gin.RouterGroup) {
	router.GET("/ws/chat/:chatID", s.SocketHandler.HandleWebSocket)
	router.GET("/ws/user/:userID", s.SocketHandler.HandleWebSocket)
	router.GET("/ws/group", s.SocketHandler.HandleWebSocket)
}

func (s *Server) Start() (err error) {
	gin.SetMode(gin.ReleaseMode)
	gin.SetMode(gin.ReleaseMode)

	docs.SwaggerInfo.Title = build.AppName
	docs.SwaggerInfo.Description = fmt.Sprintf("%s APIs", build.AppName)
	docs.SwaggerInfo.Version = build.Version
	docs.SwaggerInfo.Host = config.Cfg.HTTP.Origin
	docs.SwaggerInfo.BasePath = os.Getenv("SWAGGER_BASE")
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	app := gin.New()
	app.Use(gin.Recovery())
	log.Println(config.Cfg.HTTP)
	if len(config.Cfg.HTTP.AllowOrigins) > 0 {
		corsConfig := cors.DefaultConfig()
		corsConfig.AllowOrigins = config.Cfg.HTTP.AllowOrigins
		corsConfig.AllowCredentials = true
		corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "Authorization")
		app.Use(cors.New(corsConfig))
		log.Println(corsConfig)
	}
	app.Use(otelgin.Middleware("app-api"))
	app.Use(utils.HTTPLogger)

	store := cookie.NewStore([]byte(config.Cfg.HTTP.Secret))
	store.Options(sessions.Options{MaxAge: 60 * 60, Path: "/"})
	sessMiddleware := sessions.Sessions("app_session", store)
	app.Use(sessMiddleware)

	api := app.Group("/api/v1")
	socket := app.Group("/")

	s.Socket(socket)
	s.Routes(api)

	logger.For(nil).Info(config.Cfg.HTTP.FullAddr())
	if config.Cfg.HTTP.EnableSSL {
		return app.RunTLS(config.Cfg.HTTP.Addr(), config.Cfg.HTTP.CertFile, config.Cfg.HTTP.KeyFile)
	}
	return app.Run(config.Cfg.HTTP.Addr())
}

func abortWithStatusError(ctx *gin.Context, status int, err error) {
	if err := apperror.As(err); err != nil {
		if config.Cfg.Logger.StackTrace {
			logger.For(ctxFromGin(ctx)).Errorf("%s%s", err, err.StackTrace())
		} else {
			logger.For(ctxFromGin(ctx)).Errorf("%s", err)
		}
		if err.Code == errors.CodeDatabaseError || err.Code == errors.CodeExternalError {
			status = 500
		}
		ctx.AbortWithStatusJSON(status, dto.HTTPResp{}.FromErr(err))
		return
	}
	logger.For(ctxFromGin(ctx)).Errorf("error %v", err)
	ctx.AbortWithStatus(http.StatusInternalServerError)
}

func ctxFromGin(c *gin.Context) context.Context {
	return c.Request.Context()
}
