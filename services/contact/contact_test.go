package contact

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	db *DB
)

func TestMain(m *testing.M) {
	conn, err := gorm.Open(sqlite.Open("./testdata/contact.db"))
	if err != nil {
		log.Fatal(err)
	}

	conn.AutoMigrate(&Contact{})
	db = &DB{Conn: conn}

	os.Exit(m.Run())
}

func TestValidateContact(t *testing.T) {
	table := []struct {
		name    string
		contact Contact
		want    error
	}{
		{
			name: "All good",
			contact: Contact{
				UserID:   1,
				Fullname: "Alugbin Abiodun",
				Email:    "tolaabbey009@gmail.com",
				Phone:    "+2347033304280",
				Address:  "33, Tioya Street, Ibadan",
			},
			want: nil,
		},
		{
			name: "Invalid ID",
			contact: Contact{
				UserID:   0,
				Fullname: "Alugbin Abiodun",
				Email:    "tolaabbey009@gmail.com",
				Phone:    "+2347033304280",
				Address:  "33, Tioya Street, Ibadan",
			},
			want: errInvalidUserID,
		},
		{
			name: "Empty Full Name",
			contact: Contact{
				UserID:   1,
				Fullname: "",
				Email:    "tolaabbey009@gmail.com",
				Phone:    "+2347033304280",
				Address:  "33, Tioya Street, Ibadan",
			},
			want: errEmptyName,
		},
		{
			name: "Empty Email",
			contact: Contact{
				UserID:   1,
				Fullname: "Alugbin Abiodun",
				Email:    "",
				Phone:    "+2347033304280",
				Address:  "33, Tioya Street, Ibadan",
			},
			want: errEmptyEmail,
		},
		{
			name: "Empty Phone",
			contact: Contact{
				UserID:   1,
				Fullname: "Alugbin Abiodun",
				Email:    "tolaabbey009@gmail.com",
				Phone:    "",
				Address:  "33, Tioya Street, Ibadan",
			},
			want: errEmptyPhone,
		},
		{
			name: "Empty Address",
			contact: Contact{
				UserID:   1,
				Fullname: "Alugbin Abiodun",
				Email:    "tolaabbey009@gmail.com",
				Phone:    "+2347033304280",
				Address:  "",
			},
			want: errEmptyAddress,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.contact.validate()
			if got != tt.want {
				t.Fatalf("Test Failed\nWant:%v\nGot: %v\n", tt.want, got)
			}
		})
	}
}

func TestCreate(t *testing.T) {
	contact := Contact{
		UserID:   1,
		Fullname: "Alugbin Abiodun",
		Email:    "tolaabbey009@gmail.com",
		Phone:    "+2347033304280",
		Address:  "33, Tioya Street, Ibadan",
	}

	res, err := db.Create(contact)
	require.Nil(t, err)
	require.NotNil(t, res)
	assert.True(t, res.ID > 0)
	assert.Equal(t, res.CreatedAt, res.UpdatedAt)

	// // Test duplicate record
	contact.Fullname = "Duplicated Contact"
	res, err = db.Create(contact)
	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "contact with this email exists")

	t.Cleanup(func() {
		require.Nil(t, cleanup())
	})
}

func TestFindByUserID(t *testing.T) {
	require.Nil(t, cleanup())
	userID := uint(1)
	createForSearch(t, userID)

	c, err := db.FindByUserID(uint32(userID))
	require.Nil(t, err)
	require.NotNil(t, c)
	require.Equal(t, 2, len(c))
	for _, v := range c {
		assert.Equal(t, uint(1), v.UserID)
	}

	t.Cleanup(func() {
		require.Nil(t, cleanup())
	})
}

func TestSearch(t *testing.T) {
	userID := uint(1)
	createForSearch(t, userID)

	c, err := db.Search(uint32(userID), "Alugbin")
	require.NoError(t, err)
	require.NotNil(t, c)
	assert.Equal(t, 2, len(c))

	t.Cleanup(func() {
		require.Nil(t, cleanup())
	})
}

func TestFindByID(t *testing.T) {
	userID, fakeUserID := uint(1), uint(2)
	createForSearch(t, userID)
	res, err := db.FindByID(userID, 1)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, userID, res.UserID)

	// find by fake userid
	res, err = db.FindByID(fakeUserID, 1)
	require.Nil(t, res)
	require.NotNil(t, err)
	require.EqualError(t, err, errNotUserContact.Error())

	t.Cleanup(func() {
		require.Nil(t, cleanup())
	})
}

func TestUpdate(t *testing.T) {
	userID := uint(1)
	createForSearch(t, userID)
	res, err := db.FindByID(userID, 1)
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, userID, res.UserID)

	res.Fullname = "Updated Fullname"
	res.Phone = "08155040074"

	err = db.Update(res)
	require.NoError(t, err)
	assert.Equal(t, "08155040074", res.Phone)
	assert.Equal(t, "Updated Fullname", res.Fullname)

	t.Cleanup(func() {
		require.Nil(t, cleanup())
	})
}

func createForSearch(t *testing.T, userID uint) {
	contacts := []Contact{
		{
			UserID:   userID,
			Fullname: "Alugbin Abiodun",
			Email:    "tolaabbey009@gmail.com",
			Phone:    "+2347033304280",
			Address:  "33, Tioya Street, Ibadan",
		},
		{
			UserID:   1,
			Fullname: "Alugbin Abiodun Olutola",
			Email:    "tolaabbey001@gmail.com",
			Phone:    "+2347033304280",
			Address:  "33, Tioya Street, Ibadan",
		},
	}

	for _, contact := range contacts {
		res, err := db.Create(contact)
		require.Nil(t, err)
		require.NotNil(t, res)
		assert.True(t, res.ID > 0)
	}
}

func cleanup() error {
	return db.Conn.Exec("DELETE FROM contacts").Error
}
