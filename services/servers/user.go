package servers

import (
	"context"
	pb "grpc-contact-manager/contact"
	"grpc-contact-manager/services/user"

	"github.com/gin-gonic/gin"
)

type UserManagerGrpc struct {
	DB *user.DB
	pb.UnimplementedUserManagerServer
}

func NewUserManagerGRPC(db *user.DB) *UserManagerGrpc {
	return &UserManagerGrpc{
		DB: db,
	}
}

func (s *Server) userRoutes() {
	user := s.Router.Group("/users")
	{
		user.GET("/", s.userIndex)
		user.POST("/auth", s.authenticate)
	}
}

func (s *Server) userIndex(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "hello world",
	})
}

func (s *Server) authenticate(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Authenticating",
	})
}

func (c *UserManagerGrpc) CreateNewUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.User, error) {
	user := user.User{
		Name:     in.Name,
		Email:    in.Email,
		Password: in.Password,
	}
	u, err := c.DB.Create(user)
	if err != nil {
		return nil, err
	}
	return &pb.User{
		Id:    int32(u.ID),
		Name:  u.Name,
		Email: u.Email,
		Token: u.Token,
	}, err
}

func (c *UserManagerGrpc) Authenticate(ctx context.Context, in *pb.AuthUserRequest) (*pb.User, error) {
	return nil, nil
}

func (s *Server) NewUser() {

}
