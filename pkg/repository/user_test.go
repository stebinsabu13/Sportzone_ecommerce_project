package repository

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stebinsabu13/ecommerce-api/pkg/domain"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestFindbyEmail(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedOutput domain.User
		buildStub      func(mock sqlmock.Sqlmock)
		expectederr    error
	}{
		{
			name:  "valid email",
			input: "stebinsabu369@gmail.com",
			expectedOutput: domain.User{
				FirstName: "Stebin",
				LastName:  "Sabu",
				Email:     "stebinsabu369@gmail.com",
				Password:  "Stebin@333",
				MobileNum: "9947650091",
				Block:     false,
				Verified:  true,
			},
			buildStub: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"first_name,last_name,email,password,mobile_num,block,verified"}).AddRow("Stebin", "Sabu", "stebinsabu369@gmail.com", "Stebin@333", "9947650091", false, true)
				query := "SELECT first_name,last_name,email,password,mobile_num,block,verified from users where email=\\?"
				mock.ExpectQuery(query).WithArgs("stebinsabu369@gmail.com").WillReturnRows(rows)
			},
			expectederr: nil,
		},
		{
			name:           "not exsisting email",
			input:          "notexsisting@gmail.com",
			expectedOutput: domain.User{},
			buildStub: func(mock sqlmock.Sqlmock) {
				query := "SELECT first_name,last_name,email,password,mobile_num,block,verified from users where email=\\?"
				mock.ExpectQuery(query).WithArgs("notexsisting@gmail.com").WillReturnError(errors.New("user not found"))
			},
			expectederr: errors.New("user not found"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
			if err != nil {
				t.Fatalf("an error '%s' was not expected when initializing a mock db session", err)
			}
			userRepository := NewUserRepository(gormDB)
			tt.buildStub(mock)
			actualOutput, actualErr := userRepository.FindbyEmail(context.TODO(), tt.input)
			if tt.expectederr == nil {
				assert.NoError(t, actualErr)
			} else {
				assert.Equal(t, tt.expectederr, actualErr)
			}

			if !reflect.DeepEqual(tt.expectedOutput, actualOutput) {
				t.Errorf("got %v, but want %v", actualOutput, tt.expectedOutput)
			}

			// Check that all expectations were met
			err = mock.ExpectationsWereMet()
			if err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}
