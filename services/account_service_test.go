package services

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/melkdesousa/gamgo/dao"
	"github.com/melkdesousa/gamgo/database"
	"github.com/stretchr/testify/assert"
)

func TestAccountService(t *testing.T) {
	err := godotenv.Load("../.env.test")
	assert.NoError(t, err, "Expected no error loading .env file")
	isIntegrationTest := os.Getenv("INTEGRATION_TEST") == "true"
	if !isIntegrationTest {
		t.Skip("Skipping integration test for AccountService. Set INTEGRATION_TEST=true to run.")
	}
	accountDAO := dao.NewAccountDAO(database.GetDBConnection())
	accountService := NewAccountService(accountDAO)
	t.Run("TestGetAccount", func(t *testing.T) {
		account, err := accountService.GetAccount("john.doe@example.com", "secret")
		assert.NoError(t, err, "Expected no error when getting account")
		assert.NotNil(t, account, "Expected account to be found")
		assert.Equal(t, "john.doe@example.com", account.Email)
		assert.Equal(t, "John Doe", account.Name)
		assert.NotZero(t, account.ID)
		assert.False(t, account.CreatedAt.IsZero())
		assert.Empty(t, account.PasswordHash, "Password hash should not be returned in response")
	})
	errorTestCases := []struct {
		name         string
		email        string
		password     string
		assertionMsg string
	}{
		{"TestGetAccountInactive", "alice.johnson@example.com", "turtle", "Expected error when getting inactive account"},
		{"TestGetAccountDeleted", "jane.smith@example.com", "burble", "Expected error when getting deleted account"},
		{"TestGetAccountNotFound", "not.found@example.com", "secret", "Expected error when getting non-existent account"},
		{"TestGetAccountWrongPassword", "john.doe@example.com", "wrongpassword", "Expected error when getting account with wrong password"},
		{"TestGetAccountWithEmptyEmail", "", "secret", "Expected error when getting account with empty email"},
		{"TestGetAccountWithEmptyPassword", "john.doe@example.com", "", "Expected error when getting account with empty password"},
		{"TestGetAccountWithInvalidEmail", "invalid-email", "secret", "Expected error when getting account with invalid email format"},
		{"TestGetAccountWithInvalidPassword", "john.doe@example.com", "invalid-password", "Expected error when getting account with invalid password"},
		{"TestGetAccountWithNonExistentEmail", "non.existent@example.com", "secret", "Expected error when getting account with non-existent email"},
		{"TestGetAccountWithSQLInjection", "' OR '1'='1", "secret", "Expected error when getting account with SQL injection"},
		{"TestGetAccountWithXSSInjection", "<script>alert('XSS')</script>", "secret", "Expected error when getting account with XSS injection"},
		{"TestGetAccountWithSQLInjectionInPassword", "john.doe@example.com", "' OR '1'='1", "Expected error when getting account with SQL injection in password"},
		{"TestGetAccountWithXSSInjectionInPassword", "john.doe@example.com", "<script>alert('XSS')</script>", "Expected error when getting account with XSS injection in password"},
		{"TestGetAccountWithEmptyEmailAndPassword", "", "", "Expected error when getting account with empty email and password"},
	}
	for _, tc := range errorTestCases {
		t.Run(tc.name, func(t *testing.T) {
			account, err := accountService.GetAccount(tc.email, tc.password)
			assert.Error(t, err, tc.assertionMsg)
			assert.Nil(t, account, "Expected no account to be returned")
		})
	}
}
