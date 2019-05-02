package valet

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/midgarco/valet_manager/internal/manager"
	"github.com/midgarco/valet_manager/pkg/config"
	"github.com/midgarco/valet_manager/pkg/valet"
	"golang.org/x/crypto/bcrypt"
)

func TestUser_Create(t *testing.T) {
	type fields struct {
		FirstName string
		LastName  string
		Email     string
		Password  string
		Address   string
		City      string
		State     string
		Zipcode   string
	}
	type args struct {
		db *gorm.DB
	}

	_ = config.LoadEnv("../config", config.Option{"APP_ENV", "local"})

	conn := &manager.Connection{}
	_ = conn.DBConnection()
	defer conn.DB.Close()

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
		{name: "Create User", fields: fields{
			FirstName: "Jeff",
			LastName:  "Dupont",
			Email:     fmt.Sprintf("jeff.dupont+test-%s@gmail.com", strconv.Itoa(int(time.Now().Unix()))),
			Password:  "pass123",
			Address:   "123 Main St",
			City:      "Anycity",
			State:     "CA",
			Zipcode:   "00001",
		}, args: args{
			db: conn.DB,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &valet.User{
				FirstName: tt.fields.FirstName,
				LastName:  tt.fields.LastName,
				Email:     tt.fields.Email,
				Password:  tt.fields.Password,
				Address: valet.Address{
					Line1:   tt.fields.Address,
					City:    tt.fields.City,
					State:   tt.fields.State,
					Zipcode: tt.fields.Zipcode,
				},
			}
			if err := u.Create(tt.args.db); err != nil {
				t.Error(err)
			}

			if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(tt.fields.Password)); err != nil {
				t.Errorf("password failed")
			}
			if u.ID == 0 {
				t.Errorf("failed to create user")
			}
		})
	}
}
