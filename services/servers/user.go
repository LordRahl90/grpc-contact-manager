package servers

import (
	"context"
	"net/http"

	pb "grpc-contact-manager/contact"
	"grpc-contact-manager/services/user"

	"github.com/gin-gonic/gin"
)

type UserManagerGrpc struct {
	DB *user.DB
	pb.UnimplementedUserManagerServer
}

// CreateUserReq request struct
type CreateUserReq struct {
	Name     string `json:"name" form:"name" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
	Email    string `json:"email" form:"email" binding:"required"`
}

type AuthenticateUserReq struct {
	Email    string `json:"email" form:"email" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

// Auth view struc
type Auth struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewUserManagerGRPC(db *user.DB) *UserManagerGrpc {
	return &UserManagerGrpc{
		DB: db,
	}
}

// UserRoutes registers users routes
func (s *Server) UserRoutes() {
	user := s.Router.Group("/users")
	{
		user.GET("/", s.userIndex)
		user.POST("/", s.newUser)
		user.POST("/auth", s.authenticate)
	}
}

func (s *Server) userIndex(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "hello world",
	})
}

func (s *Server) authenticate(c *gin.Context) {
	var u AuthenticateUserReq
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newUser, err := userDB.Authenticate(u.Email, u.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User authenticated successfully",
		"data":    newUser,
	})
}

func (s *Server) newUser(c *gin.Context) {
	var u CreateUserReq
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newUser, err := userDB.Create(user.User{
		Name:     u.Name,
		Email:    u.Email,
		Password: u.Password,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err,
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "User created successfully",
		"data":    newUser,
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
	user, err := c.DB.Authenticate(in.Email, in.Password)
	if err != nil {
		return nil, err
	}
	return &pb.User{
		Id:    int32(user.ID),
		Name:  user.Name,
		Email: user.Email,
		Token: user.Token,
	}, nil
}
