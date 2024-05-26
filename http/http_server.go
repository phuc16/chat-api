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
	UserSvc    *service.UserService
	OtpSvc     *service.OtpService
	ChatSvc    *service.ChatService
	MessageSvc *service.MessageService
}

func NewServer(userSvc *service.UserService, otpSvc *service.OtpService, chatSvc *service.ChatService, messageSvc *service.MessageService) *Server {
	return &Server{UserSvc: userSvc, OtpSvc: otpSvc, ChatSvc: chatSvc, MessageSvc: messageSvc}
}

func (s *Server) Routes(router *gin.RouterGroup) {
	router.GET("/health", func(ctx *gin.Context) {
		ctx.String(200, build.Info().String())
	})
	if !config.Cfg.HTTP.IsProduction {
		router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	router.POST("/auth/register", s.CreateUser)
	router.POST("/auth/login", s.Login)
	router.GET("/auth/logout", s.Logout)

	router.GET("/user/profile", s.Authenticate, s.GetProfile)
	router.GET("/users", s.Authenticate, s.GetUserList)
	router.GET("/users/:id", s.Authenticate, s.GetUser)
	// router.PUT("/users/active", s.ActiveUser)
	router.POST("/users/check-phone", s.CheckPhoneNumber)
	router.PUT("/users/reset-password", s.ResetPassword)
	router.PUT("/users", s.Authenticate, s.UpdateUser)
	router.DELETE("/users", s.Authenticate, s.DeleteUser)
	router.POST("/users/:id/friends/request", s.Authenticate, s.SendFriendRequest)
	router.DELETE("/users/:id/friends/reject", s.Authenticate, s.RejectFriendRequest)
	router.POST("/users/:id/friends/accept", s.Authenticate, s.AcceptFriendRequest)
	router.DELETE("/users/:id/friends/remove", s.Authenticate, s.RemoveFriend)
	router.GET("/users/friends/suggest", s.Authenticate, s.SuggestFriend)

	router.POST("/chat", s.Authenticate, s.CreateChat)
	router.GET("/chat", s.Authenticate, s.GetChatList)
	router.POST("/chat/group", s.Authenticate, s.CreateGroup)
	router.PUT("/chat/group/rename", s.Authenticate, s.RenameGroup)
	router.PUT("/chat/groupAdd", s.Authenticate, s.AddToGroup)
	router.PUT("/chat/groupRemove", s.Authenticate, s.RemoveFromGroup)

	router.POST("/message", s.Authenticate, s.SendMessage)
	router.GET("/messages/chat_id", s.Authenticate, s.GetMessageListByChatId)

	router.POST("/otps/request", s.RequestOtp)
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

	api := app.Group("/api")

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
