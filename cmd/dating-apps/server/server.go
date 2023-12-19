package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"gilsaputro/dating-apps/cmd/dating-apps/config"
	auth_handler "gilsaputro/dating-apps/internal/handler/authentication"
	"gilsaputro/dating-apps/internal/handler/middleware"
	partner_handler "gilsaputro/dating-apps/internal/handler/partner"
	user_handler "gilsaputro/dating-apps/internal/handler/user"
	auth_service "gilsaputro/dating-apps/internal/service/authentication"
	partner_service "gilsaputro/dating-apps/internal/service/partner"
	user_service "gilsaputro/dating-apps/internal/service/user"
	partner_store "gilsaputro/dating-apps/internal/store/partnercache"
	user_store "gilsaputro/dating-apps/internal/store/user"
	userhist_store "gilsaputro/dating-apps/internal/store/userhistory"
	"gilsaputro/dating-apps/pkg/hash"
	"gilsaputro/dating-apps/pkg/postgres"
	"gilsaputro/dating-apps/pkg/redis"
	"gilsaputro/dating-apps/pkg/token"
	"gilsaputro/dating-apps/pkg/vault"
)

// Servcer is list configuration to run Server
type Server struct {
	cfg            config.Config
	vault          vault.VaultMethod
	hashMethod     hash.HashMethod
	tokenMethod    token.TokenMethod
	postgres       postgres.PostgresMethod
	redisMethod    redis.RedisMethod
	middleware     middleware.Middleware
	userStore      user_store.UserStoreMethod
	userService    user_service.UserServiceMethod
	userHandler    user_handler.UserHandler
	authService    auth_service.AuthenticationServiceMethod
	authHandler    auth_handler.AuthenticationHandler
	partnerStore   partner_store.PartnerCacheStoreMethod
	partnerService partner_service.PartnerServiceMethod
	partnerHandler partner_handler.PartnerHandler
	userHistStore  userhist_store.UserHistoryStoreMethod
	httpServer     *http.Server
}

// NewServer is func to create server with all configuration
func NewServer() (*Server, error) {
	s := &Server{}

	// ======== Init Dependencies Related ========
	// Load Env File
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error loading .env file: %v", err)
		return s, err
	}

	// Init Vault
	{
		token := os.Getenv("VAULT_TOKEN")
		if len(token) <= 0 {
			fmt.Print("[Got Error]-Vault Invalid VAULT_TOKEN")
			return s, fmt.Errorf("[Got Error]-Vault Invalid VAULT_TOKEN")
		}

		host := os.Getenv("VAULT_HOST")
		if len(host) <= 0 {
			fmt.Print("[Got Error]-Vault Invalid VAULT_HOST")
			return s, fmt.Errorf("[Got Error]-Vault Invalid VAULT_HOST")
		}

		vaultMethod, err := vault.NewVaultClient(token, host)
		if err != nil {
			fmt.Print("[Got Error]-Vault :", err)
		}
		s.vault = vaultMethod

		log.Println("Init-Vault")
	}

	// Get Config from yaml and replace by secret
	{
		secret, err := s.vault.GetConfig()
		if err != nil {
			fmt.Print("[Got Error]-Load Secret :", err)
		}
		cfg, err := config.GetConfig(secret)
		if err != nil {
			fmt.Print("[Got Error]-Load Config :", err)
			return s, err
		}
		s.cfg = cfg

		log.Println("LOAD-Config")
	}

	// Init Postgres
	{
		postgresMethod, err := postgres.NewPostgresClient(s.cfg.Postgres.Config)
		if err != nil {
			fmt.Print("[Got Error]-Postgres :", err)
			return s, err
		}

		s.postgres = postgresMethod

		log.Println("Init-Postgres")
	}

	// Init Redis
	{
		redisMethod := redis.NewRedisClient(redis.RedisConfig{
			Host:     s.cfg.Redis.Host,
			Port:     s.cfg.Redis.Port,
			Password: s.cfg.Redis.Password,
		})
		s.redisMethod = redisMethod
		log.Println("Init-Redis")
	}

	// Init Hash Package
	{
		hashMethod := hash.NewHashMethod(s.cfg.Hash.Cost)
		s.hashMethod = hashMethod
		log.Println("Init-Hash Package")
	}

	// Init Token Package
	{
		tokenMethod := token.NewTokenMethod(s.cfg.Token.Secret, s.cfg.Token.ExpInHour)
		s.tokenMethod = tokenMethod
		log.Println("Init-Token Package")
	}

	// ======== Init Dependencies Store ========
	// Init User Store
	{
		userStore := user_store.NewUserStore(s.postgres)
		s.userStore = userStore
		log.Println("Init-User Store")
	}

	{
		userHistStore := userhist_store.NewUserHistoryStore(s.postgres)
		s.userHistStore = userHistStore
		log.Println("Init-User History Store")
	}

	{
		partnerStore := partner_store.NewPartnerCacheStore(s.redisMethod)
		s.partnerStore = partnerStore
		log.Println("Init-Partner Cache Store")
	}

	// ======== Init Dependencies Service ========
	// Init User Service
	{
		userService := user_service.NewUserService(s.userStore, s.hashMethod)
		s.userService = userService
		log.Println("Init-User Service")
	}

	{
		authService := auth_service.NewAuthenticationService(s.userStore, s.tokenMethod, s.hashMethod)
		s.authService = authService
		log.Println("Init-Auth Service")
	}

	{
		partnerService := partner_service.NewPartnerService(s.userStore, s.userHistStore, s.partnerStore, s.cfg.MaxCounter)
		s.partnerService = partnerService
		log.Println("Init-Partner Service")
	}

	// ======== Init Dependencies Handler ========
	// Init Middleware
	{
		midlewareService := middleware.NewMiddleware(s.tokenMethod)
		s.middleware = midlewareService
		log.Println("Init-Middleware")
	}

	// Init User Handler
	{
		var opts []user_handler.Option
		opts = append(opts, user_handler.WithTimeoutOptions(s.cfg.UserHandler.TimeoutInSec))
		userHandler := user_handler.NewUserHandler(s.userService, opts...)
		s.userHandler = *userHandler
		log.Println("Init-User Handler")
	}

	// Init Auth Handler
	{
		var opts []auth_handler.Option
		opts = append(opts, auth_handler.WithTimeoutOptions(s.cfg.AuthHandler.TimeoutInSec))
		authHandler := auth_handler.NewAuthenticationHandler(s.authService, opts...)
		s.authHandler = *authHandler
		log.Println("Init-Auth Handler")
	}

	// Init Partner Handler
	{
		var opts []partner_handler.Option
		opts = append(opts, partner_handler.WithTimeoutOptions(s.cfg.PartnerHandler.TimeoutInSec))
		partnerHandler := partner_handler.NewPartnerHandler(s.partnerService, opts...)
		s.partnerHandler = *partnerHandler
		log.Println("Init-Partner Handler")
	}

	// Init Router
	{
		r := mux.NewRouter()
		// Init Guest Path
		r.HandleFunc("/v1/login", s.authHandler.LoginUserHandler).Methods("POST")
		r.HandleFunc("/v1/register", s.authHandler.RegisterUserHandler).Methods("POST")

		// Init User Path
		r.HandleFunc("/v1/user", s.middleware.MiddlewareVerifyToken(s.userHandler.DeleteUserHandler)).Methods("DELETE")
		r.HandleFunc("/v1/user", s.middleware.MiddlewareVerifyToken(s.userHandler.EditUserHandler)).Methods("PUT")

		// Init Partner Partner Path
		r.HandleFunc("/v1/partner", s.middleware.MiddlewareVerifyToken(s.partnerHandler.CurrentPartnerHandler)).Methods("GET")
		r.HandleFunc("/v1/partner/history", s.middleware.MiddlewareVerifyToken(s.partnerHandler.LikedHistoryHandler)).Methods("GET")
		r.HandleFunc("/v1/partner/pass", s.middleware.MiddlewareVerifyToken(s.partnerHandler.PassPartnerHandler)).Methods("POST")
		r.HandleFunc("/v1/partner/like", s.middleware.MiddlewareVerifyToken(s.partnerHandler.LikePartnerHandler)).Methods("POST")

		port := ":" + s.cfg.Port
		log.Println("running on port ", port)

		server := &http.Server{
			Addr:    port,
			Handler: r,
		}

		s.httpServer = server
	}
	return s, nil
}

func (s *Server) Start() int {
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil {
			fmt.Println(err)
		}
	}()

	// Wait for a signal to shut down the application
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// Create a context with a timeout to allow the server to cleanly shut down
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	s.httpServer.Shutdown(ctx)
	log.Println("complete, shutting down.")
	return 0
}

// Run is func to create server and invoke Start()
func Run() int {
	s, err := NewServer()
	if err != nil {
		return 1
	}

	return s.Start()
}
