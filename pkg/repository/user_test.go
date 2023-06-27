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
				query := `SELECT \* from users where email=\$1 and deleted_at IS NULL`
				mock.ExpectQuery(query).WithArgs("notexsisting@gmail.com").WillReturnError(errors.New("invalid email"))
			},
			expectederr: errors.New("invalid email"),
		},
		// {
		// 	name:  "valid email",
		// 	input: "stebinsabu369@gmail.com",
		// 	expectedOutput: utils.ResponseUsers{
		// 		FirstName: "Stebin",
		// 		LastName:  "Sabu",
		// 		Email:     "stebinsabu369@gmail.com",
		// 		Password:  "Stebin@333",
		// 		MobileNum: "9947650091",
		// 	},
		// 	buildStub: func(mock sqlmock.Sqlmock) {
		// 		rows := sqlmock.NewRows([]string{"first_name", "last_name", "email", "password", "mobile_num", "block", "verified"}).AddRow("Stebin", "Sabu", "stebinsabu369@gmail.com", "Stebin@333", "9947650091", false, true)
		// 		query := `SELECT \* from users where email=\$1 and deleted_at IS NULL`
		// 		mock.ExpectQuery(query).WithArgs("stebinsabu369@gmail.com").WillReturnRows(rows)
		// 	},
		// 	expectederr: nil,
		// },
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
		buildstub      func(mock sqlmock.Sqlmock)
		expectedOutput string
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
			buildstub: func(mock sqlmock.Sqlmock) {
				// mockTime := time.Date(2023, 6, 23, 17, 6, 47, 201, time.UTC)
				// current := time.Now()
				// trunctated := current.Truncate(time.Millisecond)
				// query := `insert into users\(created_at,updated_at,first_name,last_name,email,mobile_num,password,referal_code\)values\(\'.+\'\,\'.+\'\,\'Stebin\'\,\'Sabu\'\,\'stebinsabu369@gmail.com\'\,\'9947650091\'\,\'Stebin@333\'\,\'jdj43\'\) returning id`
				mock.ExpectBegin()
				// row := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectExec(`insert into users(created_at,updated_at,first_name,last_name,email,mobile_num,password,referal_code)values($1,$2,$3,$4,$5,$6,$7,$8) returning id`).
					WithArgs("2023-06-23 18:23:05.398", "2023-06-23 18:23:05.398", "Stebin", "Sabu", "stebinsabu369@gmail.com", "9947650091", "Stebin@333", "jdj43").
					WillReturnResult(sqlmock.NewResult(1, 1))
				// mock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("insert into carts(user_id)values($1)").WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedOutput: "9947650091",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.buildstub(mockc)
			actualOutput, actualerr := userRepository.SignUpUser(context.TODO(), tt.input)
			if actualerr != nil {
				t.Errorf("An error occurred signing up the user: %v", err)
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
