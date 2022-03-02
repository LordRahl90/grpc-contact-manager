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

var (
	userDB    *user.DB
	contactDB *contact.DB
)

type Database interface {
	// Create(interface{}) (gorm.Model, error)
	// New(*gorm.DB) (*Database, error)
	Migrate() error
}

// Server creates a struct to house the server elements.
type Server struct {
	Conn   *gorm.DB
	Router *gin.Engine
}

// New initialize a new server object
func New(db *gorm.DB) (*Server, error) {
	router := gin.Default()
	return &Server{
		Conn:   db,
		Router: router,
	}, nil
}

func (s *Server) StartHttp(ctx context.Context, port string) (*http.Server, error) {
	if s == nil {
		return nil, errServerEmpty
	}
	if s.Conn == nil {
		return nil, errDatabaseNotSet
	}
	if s.Router == nil {
		return nil, errRouterNotSet
	}
	// initialize all db repositories
	// s.userRoutes()
	s.setupModels()

	return &http.Server{
		Addr:    port,
		Handler: s.Router,
	}, nil
}

// setupModels sets up server models
func (s *Server) setupModels() error {
	u, err := user.New(s.Conn)
	if err != nil {
		return err
	}
	c, err := contact.New(s.Conn)
	if err != nil {
		return err
	}
	userDB = u
	contactDB = c

	if err := userDB.Migrate(); err != nil {
		return err
	}

	return contactDB.Migrate()
}

func (s *Server) StartUserGRPC(ctx context.Context) (*grpc.Server, error) {
	userDB := &user.DB{Conn: s.Conn}
	if err := userDB.Migrate(); err != nil {
		return nil, err
	}
	userGrpcServer := NewUserManagerGRPC(userDB)
	gServer := grpc.NewServer()
	pb.RegisterUserManagerServer(gServer, userGrpcServer)
	return gServer, nil
}
