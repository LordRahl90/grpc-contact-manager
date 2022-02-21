package servers

import (
	"context"
	"errors"
	"net/http"

	pb "grpc-contact-manager/contact"
	"grpc-contact-manager/services/contact"
	"grpc-contact-manager/services/user"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

var (
	errServerEmpty    = errors.New("server not initialized")
	errDatabaseNotSet = errors.New("database not set")
	errRouterNotSet   = errors.New("router not set")
)

// Server creates a struct to house the server elements.
type Server struct {
	DB     *gorm.DB
	Router *gin.Engine
}

// New initialize a new server object
func New(db *gorm.DB) (*Server, error) {
	router := gin.Default()
	return &Server{
		DB:     db,
		Router: router,
	}, nil
}

func (s *Server) StartHttp(ctx context.Context, port string) (*http.Server, error) {
	if s == nil {
		return nil, errServerEmpty
	}
	if s.DB == nil {
		return nil, errDatabaseNotSet
	}
	if s.Router == nil {
		return nil, errRouterNotSet
	}
	// initialize all db repositories
	s.userRoutes()
	if err := user.Migrate(s.DB); err != nil {
		return nil, err
	}
	if err := contact.Migrate(s.DB); err != nil {
		return nil, err
	}

	return &http.Server{
		Addr:    port,
		Handler: s.Router,
	}, nil
}

func (s *Server) StartUserGRPC(ctx context.Context) (*grpc.Server, error) {
	userGrpcServer := NewUserManagerGRPC(&user.DB{Conn: s.DB})
	gServer := grpc.NewServer()
	pb.RegisterUserManagerServer(gServer, userGrpcServer)
	return gServer, nil
}
