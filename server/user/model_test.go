package user

import (
	"database/sql/driver"
	"log"
	"os"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db   DB
	mock sqlmock.Sqlmock
)

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
	mock = sqlMock

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
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users" ("created_at","updated_at","deleted_at","name","email","password","token") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "id"`)).
		WithArgs(AnyTime{}, AnyTime{}, nil, "Alugbin LordRahl", "tolaabbey009@gmail.com", "password", "").
		WillReturnRows(sqlmock.NewRows([]string{"ID"}).AddRow(strconv.Itoa(1)))
	mock.ExpectCommit()

	user := User{
		Name:     "Alugbin LordRahl",
		Email:    "tolaabbey009@gmail.com",
		Password: "password",
	}

	res, err := db.Create(user)
	require.Nil(t, err)
	require.NotNil(t, res)
}

type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func TestCreateWithNoName(t *testing.T) {
	dbase, _, err := sqlmock.New()
	require.Nil(t, err)
	conn, err := gorm.Open(postgres.New(postgres.Config{
		Conn: dbase,
	}))
	require.NoError(t, err)
	db := DB{Conn: conn}

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
	dbase, _, err := sqlmock.New()
	require.Nil(t, err)
	conn, err := gorm.Open(postgres.New(postgres.Config{
		Conn: dbase,
	}))
	require.NoError(t, err)
	db := DB{Conn: conn}

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

}

func TestGenerateToken(t *testing.T) {

}

func TestValidateToken(t *testing.T) {}
