// user_test.go
package database

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/skeleton1231/gotal/internal/apiserver/store/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupMockDB() (*gorm.DB, sqlmock.Sqlmock, error) {
	mockDB, mock, err := sqlmock.New() // 创建新的模拟数据库
	if err != nil {
		return nil, nil, err
	}

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      mockDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}

	return gormDB, mock, nil
}

func TestCreateUser(t *testing.T) {
	db, mock, err := setupMockDB()
	assert.NoError(t, err)

	u := newUsers(&datastore{db})

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `users`").
		WithArgs(
			sqlmock.AnyArg(),      // ExtendShadow
			sqlmock.AnyArg(),      // CreatedAt
			sqlmock.AnyArg(),      // UpdatedAt
			sqlmock.AnyArg(),      // DeletedAt
			sqlmock.AnyArg(),      // Status
			"John Doe",            // Name
			"johndoe@example.com", // Email
			sqlmock.AnyArg(),      // EmailVerifiedAt
			sqlmock.AnyArg(),      // Password
			sqlmock.AnyArg(),      // RememberToken
			sqlmock.AnyArg(),      // StripeID
			sqlmock.AnyArg(),      // DiscordID
			sqlmock.AnyArg(),      // PMType
			sqlmock.AnyArg(),      // PMLastFour
			sqlmock.AnyArg(),      // TrialEndsAt
			sqlmock.AnyArg(),      // TotalCredits
			// ... other fields if there are any
		).WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	err = u.Create(context.Background(), &model.User{Name: "John Doe", Email: "johndoe@example.com"}, model.CreateOptions{})
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// 更多测试函数...
