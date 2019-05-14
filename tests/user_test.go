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
		FirstName   string
		LastName    string
		Email       string
		Password    string
		Address     string
		City        string
		State       string
		Zipcode     string
		HomePhone   string
		WorkPhone   string
		MobilePhone string
	}
	type args struct {
		db *gorm.DB
	}

	_ = config.LoadEnv("../config", config.Option{"APP_ENV", "local"})

	conn := &manager.Connection{}
	err := conn.DBConnection()
	if err != nil {
		t.Error(err)
		return
	}
	defer conn.DB.Close()

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
		{name: "Create User", fields: fields{
			FirstName:   "Jeff",
			LastName:    "Dupont",
			Email:       fmt.Sprintf("jeff.dupont+test-%s@gmail.com", strconv.Itoa(int(time.Now().Unix()))),
			Password:    "pass123",
			Address:     "123 Main St",
			City:        "Anycity",
			State:       "CA",
			Zipcode:     "00001",
			HomePhone:   "222 123-4567",
			WorkPhone:   "333 456-7890",
			MobilePhone: "444 567-8901",
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
				PhoneNumbers: []valet.Phone{
					valet.Phone{Type: "home", Number: tt.fields.HomePhone},
					valet.Phone{Type: "work", Number: tt.fields.WorkPhone},
					valet.Phone{Type: "mobile", Number: tt.fields.MobilePhone},
				},
			}
			if err := u.Create(tt.args.db); err != nil {
				t.Error(err)
				return
			}

			if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(tt.fields.Password)); err != nil {
				t.Errorf("password failed")
				return
			}
			if u.ID == 0 {
				t.Errorf("failed to create user")
				return
			}
		})
	}
}
