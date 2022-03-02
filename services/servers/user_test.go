package servers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	pb "grpc-contact-manager/contact"
	"grpc-contact-manager/services/user"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestLoadUser(t *testing.T) {
	ctx := context.Background()
	s, err := server.StartHttp(ctx, ":2500")
	require.NoError(t, err)
	require.NotNil(t, s)

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/users/", nil)
	require.NoError(t, err)
	s.Handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, `{"message":"hello world"}`, w.Body.String())
	t.Cleanup(func() {
		require.NoError(t, cleanup(server.Conn))
	})
}

func TestCreateUser(t *testing.T) {
	ctx := context.Background()
	s, err := server.StartHttp(ctx, ":2500")
	require.NoError(t, err)
	require.NotNil(t, s)
	require.NotNil(t, userDB)

	payload := `{
		"name":"Alugin Abiodun",
		"email":"tolaabbey009@gmail.com",
		"password":"password"
	}`

	req, err := http.NewRequest("POST", "/users/", strings.NewReader(payload))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	s.Handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	resp := make(map[string]interface{})

	err = json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	require.NotNil(t, resp)
	data, ok := resp["data"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, float64(1), data["ID"].(float64))
	assert.Empty(t, data["password"].(string))

	t.Cleanup(func() {
		require.Nil(t, cleanup(userDB.Conn))
	})
}

func TestFailedCreateUser(t *testing.T) {
	ctx := context.Background()
	s, err := server.StartHttp(ctx, ":2500")
	require.NoError(t, err)
	require.NotNil(t, s)
	require.NotNil(t, userDB)

	payload := `{
		"names":"Alugin Abiodun",
		"email":"tolaabbey009@gmail.com",
		"password":"password"
	}`

	req, err := http.NewRequest("POST", "/users/", strings.NewReader(payload))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	s.Handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	t.Cleanup(func() {
		require.Nil(t, cleanup(userDB.Conn))
	})
}

func TestCreateUserWithEmptyName(t *testing.T) {
	ctx := context.Background()
	s, err := server.StartHttp(ctx, ":2500")
	require.NoError(t, err)
	require.NotNil(t, s)
	require.NotNil(t, userDB)

	payload := `{
		"name":"",
		"email":"tolaabbey009@gmail.com",
		"password":"password"
	}`

	req, err := http.NewRequest("POST", "/users/", strings.NewReader(payload))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	s.Handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	t.Cleanup(func() {
		require.Nil(t, cleanup(userDB.Conn))
	})
}

func TestAuthenticateUser(t *testing.T) {
	ctx := context.Background()

	s, err := server.StartHttp(ctx, ":2500")
	require.NoError(t, err)
	require.NotNil(t, s)
	require.NotNil(t, userDB)

	// create user to authenticate against
	u, err := userDB.Create(user.User{
		Name:     "Alugbin LordRahl",
		Email:    "tolaabbey009@gmail.com",
		Password: "password",
	})
	require.NoError(t, err)
	require.NotNil(t, u)
	assert.NotEqual(t, uint(0), u.ID)

	payload := `{
		"email":"tolaabbey009@gmail.com",
		"password":"password"
	}`

	req, err := http.NewRequest("POST", "/users/auth", strings.NewReader(payload))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	s.Handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := make(map[string]interface{})

	err = json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	require.NotNil(t, resp)

	data, ok := resp["data"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, float64(u.ID), data["ID"].(float64))
	assert.Empty(t, data["password"].(string))
	assert.NotEmpty(t, data["token"].(string))
	assert.True(t, resp["success"].(bool))

	t.Cleanup(func() {
		require.Nil(t, cleanup(userDB.Conn))
	})
}

func TestAuthenticateUserWithBadRequest(t *testing.T) {
	ctx := context.Background()

	s, err := server.StartHttp(ctx, ":2500")
	require.NoError(t, err)
	require.NotNil(t, s)
	require.NotNil(t, userDB)

	// create user to authenticate against
	u, err := userDB.Create(user.User{
		Name:     "Alugbin LordRahl",
		Email:    "tolaabbey009@gmail.com",
		Password: "password",
	})
	require.NoError(t, err)
	require.NotNil(t, u)
	assert.NotEqual(t, uint(0), u.ID)

	payload := `{
		"emails":"tolaabbey009@gmail.com",
		"password":"password"
	}`

	req, err := http.NewRequest("POST", "/users/auth", strings.NewReader(payload))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	s.Handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	t.Cleanup(func() {
		require.Nil(t, cleanup(userDB.Conn))
	})
}

func TestGRPCCreateUser(t *testing.T) {
	ctx := context.Background()
	in := pb.CreateUserRequest{
		Name:     "Alugbin Abiodun",
		Email:    "tolaabbey009@gmail.com",
		Password: "password",
	}
	res, err := usergrpc.CreateNewUser(ctx, &in)
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, int32(1), res.Id)
	assert.Equal(t, in.Email, res.Email)

	t.Cleanup(func() {
		require.Nil(t, cleanup(usergrpc.DB.Conn))
	})
}

func cleanup(db *gorm.DB) error {
	return db.Exec("DELETE FROM users").Error
}
