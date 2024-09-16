package api

import (
	"fmt"

	db "github.com/Srinath-exe/simplebank/db/sqlc"
	"github.com/Srinath-exe/simplebank/token"
	"github.com/Srinath-exe/simplebank/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
	config     util.Config
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{store: store, tokenMaker: tokenMaker, config: config}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()

	return server, nil
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func (server *Server) setupRouter() {
	router := gin.Default()
	router.POST("/users", server.createUser)

	router.POST("users/login", server.loginUser)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.GET("/accounts", server.getAccountsList)

	authRoutes.POST("/transfers", server.createTransfer)
	authRoutes.POST("/transfers/account", server.listTransfersFromAccountId)
	authRoutes.GET("/users/:username", server.getUser)

	// search routes
	authRoutes.POST("/accounts/search", server.searchAccounts)
	authRoutes.POST("/transfers/search", server.searchTransfers)
	authRoutes.POST("/entries/search", server.searchEntries)

	server.router = router
}
