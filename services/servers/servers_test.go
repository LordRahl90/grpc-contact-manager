package servers

import (
	"log"
	"os"
	"testing"

	"grpc-contact-manager/services/user"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	server   *Server
	usergrpc *UserManagerGrpc
)

func TestMain(m *testing.M) {
	conn, err := gorm.Open(sqlite.Open("./testdata/contact.db"))
	if err != nil {
		log.Fatal(err)
	}
	s, err := New(conn)
	if err != nil {
		log.Fatal(err)
	}
	server = s
	usergrpc = &UserManagerGrpc{
		DB: &user.DB{Conn: conn},
	}
	server.UserRoutes()
	os.Exit(m.Run())
}
