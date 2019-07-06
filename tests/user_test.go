package valet

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/midgarco/valet_manager/pkg/pagination"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/midgarco/env"
	"github.com/midgarco/valet_manager/internal/manager"
	"github.com/midgarco/valet_manager/pkg/valet"
	"golang.org/x/crypto/bcrypt"
)

var conn *manager.Connection

func setupTestCase(t *testing.T) func(*testing.T) {
	_ = env.Load("VALET_MGR", "../config", env.Option{"APP_ENV", "local"})

	conn = &manager.Connection{}
	err := conn.DBConnection()
	if err != nil {
		t.Error(err)
	}

	return func(t *testing.T) {
		conn.DB.Close()
	}
}

func getUser(v int64) *valet.User {
	seed := rand.NewSource(v)
	r := rand.New(seed)
	return &valet.User{
		FirstName: "John",
		LastName:  "Example",
		Email:     fmt.Sprintf("john.example+test-%s@gmail.com", strconv.Itoa(r.Intn(1000))),
		Password:  "pass123",
		Address: valet.Address{
			Line1:   "123 Main St",
			City:    "Anycity",
			State:   "CA",
			Zipcode: "00001",
		},
		PhoneNumbers: []valet.Phone{
			valet.Phone{Type: "home", Number: "222 123-4567"},
			valet.Phone{Type: "work", Number: "333 456-7890"},
			valet.Phone{Type: "mobile", Number: "444 567-8901"},
		},
	}
}

func TestUser_Create(t *testing.T) {
	teardown := setupTestCase(t)
	defer teardown(t)

	u := getUser(time.Now().Unix())
	if err := u.Create(conn.DB); err != nil {
		t.Error(err)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte("pass123")); err != nil {
		t.Errorf("password failed")
		return
	}
	if u.ID == 0 {
		t.Errorf("failed to create user")
		return
	}
}

func TestUser_Update(t *testing.T) {
	teardown := setupTestCase(t)
	defer teardown(t)

	u := getUser(time.Now().Unix())
	if err := u.Create(conn.DB); err != nil {
		t.Error(err)
		return
	}

	// get created user
	uu, err := valet.FindUser(conn.DB, int(u.ID))
	if err != nil {
		t.Error(err)
		return
	}

	// update field
	uu.LastName = "Doe"
	uu.Save(conn.DB)

	uuu, err := valet.FindUser(conn.DB, int(u.ID))
	if err != nil {
		t.Error(err)
		return
	}

	if uu.LastName != uuu.LastName {
		t.Errorf("Update failed: want %s, got %s", uu.LastName, uuu.LastName)
		return
	}

	if err := valet.RemoveTestUser(conn.DB, int(u.ID)); err != nil {
		t.Error(err)
		return
	}
}

func TestUser_GetUsers(t *testing.T) {
	teardown := setupTestCase(t)
	defer teardown(t)

	tc, err := valet.UserCount(conn.DB)
	if err != nil {
		t.Error(err)
	}

	pg := pagination.Paging{
		Limit:  5,
		Offset: 0,
		OrderBy: []pagination.Order{
			pagination.Order{
				Field:     "email",
				Direction: "DESC",
			},
		},
	}
	count := 0

	for {
		users, err := valet.FindUsers(conn.DB, pg)
		if err != nil {
			t.Error(err)
		}
		if len(users) == 0 {
			break
		}
		count += len(users)
		pg.Offset = count
	}
	if tc != count {
		t.Errorf("user count got %d, want %d", count, tc)
	}
}
