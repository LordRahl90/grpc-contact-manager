package user

import (
	"database/sql/driver"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db     DB
	dbMock sqlmock.Sqlmock
)

type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

type AnyPassword struct{}

func (a AnyPassword) Match(v driver.Value) bool {
	s := v.(string)
	if len(s) < 60 {
		return false
	}
	if !strings.HasPrefix(s, "$") {
		return false
	}
	return true
}

func TestMain(m *testing.M) {
	d, sqlMock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	conn, err := gorm.Open(postgres.New(postgres.Config{
		Conn: d,
	}))
	if err != nil {
		log.Fatal(err)
	}

	db = DB{Conn: conn}
	dbMock = sqlMock

	os.Exit(m.Run())
}

func TestValidate(t *testing.T) {
	table := []struct {
		name string
		user User
		want error
	}{
		{
			name: "Full info",
			user: User{
				Name:     "Alugbin LordRahl",
				Email:    "tolaabbey009@gmail.com",
				Password: "password",
			},
			want: nil,
		},
		{
			name: "No Name",
			user: User{
				Email:    "tolaabbey009@gmail.com",
				Password: "password",
			},
			want: errNoName,
		},
		{
			name: "No Email",
			user: User{
				Name:     "Alugbin LordRahl",
				Password: "password",
			},
			want: errNoEmail,
		},
		{
			name: "No Password",
			user: User{
				Name:  "Alugbin LordRahl",
				Email: "tolaabbey009@gmail.com",
			},
			want: errNoPassword,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.user.validate()
			if got != tt.want {
				t.Fatalf("Expected: %v\t Got: %v\n", tt.want, got)
			}
		})
	}
}

func TestCreate(t *testing.T) {
	dbMock.ExpectBegin()
	dbMock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users" ("created_at","updated_at","deleted_at","name","email","password","token") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "id"`)).
		WithArgs(AnyTime{}, AnyTime{}, nil, "Alugbin LordRahl", "tolaabbey009@gmail.com", AnyPassword{}, "").
		WillReturnRows(sqlmock.NewRows([]string{"ID"}).AddRow(strconv.Itoa(1)))
	dbMock.ExpectCommit()

	user := User{
		Name:     "Alugbin LordRahl",
		Email:    "tolaabbey009@gmail.com",
		Password: "password",
	}

	res, err := db.Create(user)
	require.Nil(t, err)
	require.NotNil(t, res)
	assert.Equal(t, res.Password, "")
	assert.Equal(t, res.Email, "tolaabbey009@gmail.com")
	assert.True(t, res.CreatedAt.Before(time.Now()))
}

func TestCreateWithNoName(t *testing.T) {
	user := User{
		Name:     "",
		Email:    "tolaabbey009@gmail.com",
		Password: "password",
	}

	res, err := db.Create(user)
	require.Nil(t, res)
	require.NotNil(t, err)
	require.EqualError(t, err, "name must be provided")
}

func TestCreateWithNoEmail(t *testing.T) {
	user := User{
		Name:     "Alugbin Abiodun",
		Email:    "",
		Password: "password",
	}

	res, err := db.Create(user)
	require.Nil(t, res)
	require.NotNil(t, err)
	require.EqualError(t, err, "email must be provided")
}

func TestAuthenticate(t *testing.T) {
	fakePassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	require.Nil(t, err)
	dbMock.ExpectBegin()
	dbMock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users" ("created_at","updated_at","deleted_at","name","email","password","token") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "id"`)).
		WithArgs(AnyTime{}, AnyTime{}, nil, "Alugbin LordRahl", "tolaabbey009@gmail.com", AnyPassword{}, "").
		WillReturnRows(sqlmock.NewRows([]string{"ID"}).
			AddRow(strconv.Itoa(1)))
	dbMock.ExpectCommit()
	dbMock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT 1`)).
		WithArgs("tolaabbey009@gmail.com").
		WillReturnRows(sqlmock.NewRows([]string{"ID", "created_at", "updated_at", "deleted_at", "email", "password", "token"}).
			AddRow(uint(1), time.Now(), time.Now(), nil, "tolaabbey009@gmail.com", fakePassword, ""))

	user := User{
		Name:     "Alugbin LordRahl",
		Email:    "tolaabbey009@gmail.com",
		Password: "password",
	}

	res, err := db.Create(user)
	require.Nil(t, err)
	require.NotNil(t, res)

	authUser, err := db.Authenticate(user.Email, user.Password)
	require.Nil(t, err)
	require.NotNil(t, authUser)
	assert.NotEmpty(t, authUser.Token)
	assert.Empty(t, authUser.Password)
}

func TestGenerateToken(t *testing.T) {
	userID := 1
	token, err := generateToken(uint32(userID))
	require.Nil(t, err)
	assert.NotEmpty(t, token)
}

func TestValidateToken(t *testing.T) {
	userID := 1
	tokenString, err := generateToken(uint32(userID))
	require.Nil(t, err)
	assert.NotEmpty(t, tokenString)
}

func TestValidateInvalidToken(t *testing.T) {
	authUserID, err := validateToken("hello one two three")
	require.NotNil(t, err)
	require.Equal(t, uint32(0), authUserID)
}
