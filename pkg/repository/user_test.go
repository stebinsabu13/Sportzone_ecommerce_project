package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stebinsabu13/ecommerce-api/pkg/utils"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestFindbyEmail(t *testing.T) {
	db, mockc, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when initializing a mock db session", err)
	}
	userRepository := NewUserRepository(gormDB)
	tests := []struct {
		name           string
		input          string
		expectedOutput utils.ResponseUsers
		buildStub      func(mock sqlmock.Sqlmock)
		expectederr    error
	}{
		{
			name:           "not exsisting email",
			input:          "notexsisting@gmail.com",
			expectedOutput: utils.ResponseUsers{},
			buildStub: func(mock sqlmock.Sqlmock) {
				query := `SELECT \* from users where email=\$1`
				mock.ExpectQuery(query).WithArgs("notexsisting@gmail.com").WillReturnError(errors.New("invalid email"))
			},
			expectederr: errors.New("invalid email"),
		},
		{
			name:  "valid email",
			input: "stebinsabu369@gmail.com",
			expectedOutput: utils.ResponseUsers{
				ID:          1,
				FirstName:   "Stebin",
				LastName:    "Sabu",
				Email:       "stebinsabu369@gmail.com",
				Password:    "Stebin@333",
				MobileNum:   "9947650091",
				Block:       false,
				Verified:    true,
				ReferalCode: "jsjil",
			},
			buildStub: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "password", "mobile_num", "block", "verified", "referal_code"}).AddRow(1, "Stebin", "Sabu", "stebinsabu369@gmail.com", "Stebin@333", "9947650091", false, true, "jsjil")
				query := `SELECT \* from users where email=\$1`
				mock.ExpectQuery(query).WithArgs("stebinsabu369@gmail.com").WillReturnRows(rows)
			},
			expectederr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.buildStub(mockc)
			actualOutput, actualErr := userRepository.FindbyEmail(context.TODO(), tt.input)
			if tt.expectederr == nil {
				assert.NoError(t, actualErr)
			} else {
				assert.Equal(t, tt.expectederr, actualErr)
			}

			assert.Equal(t, tt.expectedOutput, actualOutput)

			err = mockc.ExpectationsWereMet()
			if err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestSignUpUser(t *testing.T) {
	db, mockc, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when initializing a mock db session", err)
	}
	userRepository := NewUserRepository(gormDB)
	tests := []struct {
		name           string
		input          utils.BodySignUpuser
		buildstub      func(mock sqlmock.Sqlmock, user utils.BodySignUpuser)
		expectedOutput string
		expectedErr    error
	}{
		{
			name: "SignUp user",
			input: utils.BodySignUpuser{
				FirstName:   "Stebin",
				LastName:    "Sabu",
				Email:       "stebinsabu369@gmail.com",
				MobileNum:   "9947650091",
				Password:    "Stebin@333",
				ReferalCode: "jdj43",
			},
			buildstub: func(mock sqlmock.Sqlmock, user utils.BodySignUpuser) {
				mock.ExpectBegin()
				mock.ExpectQuery("insert into users").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectExec("insert into carts").
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedOutput: "9947650091",
			expectedErr:    nil,
		},
		{
			name: "failed to signupuser",
			input: utils.BodySignUpuser{
				FirstName:   "Stebin",
				LastName:    "Sabu",
				Email:       "stebinsabu369@gmail.com",
				MobileNum:   "9947650091",
				Password:    "Stebin@333",
				ReferalCode: "jdj43",
			},
			buildstub: func(mock sqlmock.Sqlmock, user utils.BodySignUpuser) {
				mock.ExpectBegin()
				mock.ExpectQuery("insert into users").
					WillReturnError(errors.New("unique constraint violation"))
				mock.ExpectRollback()
			},
			expectedOutput: "9947650091",
			expectedErr:    errors.New("unique constraint violation"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.buildstub(mockc, tt.input)
			actualOutput, actualerr := userRepository.SignUpUser(context.TODO(), tt.input)
			if tt.expectedErr == nil {
				assert.NoError(t, actualerr)
			} else {
				assert.Equal(t, tt.expectedErr, actualerr)
			}

			assert.Equal(t, tt.expectedOutput, actualOutput)

			// Check that all expectations were met
			err = mockc.ExpectationsWereMet()
			if err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}
