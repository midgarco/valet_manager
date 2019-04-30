package valet

import (
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/midgarco/valet_manager/internal/manager"
	"github.com/midgarco/valet_manager/pkg/config"
	"github.com/midgarco/valet_manager/pkg/valet"
)

func TestUser_Create(t *testing.T) {
	type fields struct {
		FirstName string
		LastName  string
		Email     string
		Password  string
	}
	type args struct {
		db *gorm.DB
	}

	_ = config.LoadEnv("../config", config.Option{"APP_ENV", "local"})

	conn := &manager.Connection{}
	_ = conn.MySQLConnection()
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
			Email:     "jeff.dupont@gmail.com",
			Password:  "pass123",
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
			}
			u.SetPassword(tt.fields.Password)
			if err := u.Create(tt.args.db); err != nil {
				t.Error(err)
			}

			if u.ID == 0 {
				t.Errorf("failed to create user")
			}
		})
	}
}
