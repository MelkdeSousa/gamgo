package services

import (
	"errors"

	"github.com/melkdesousa/gamgo/dao"
	"github.com/melkdesousa/gamgo/dao/models"
	"golang.org/x/crypto/bcrypt"
)

type AccountService struct {
	accountDAO *dao.AccountDAO
}

func NewAccountService(accountDAO *dao.AccountDAO) *AccountService {
	return &AccountService{
		accountDAO: accountDAO,
	}
}

func (s *AccountService) GetAccount(email, password string) (*models.Account, error) {
	if password == "" {
		return nil, errors.New("password cannot be empty")
	}
	account, err := s.accountDAO.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	isValid, err := ComparePasswords(account.PasswordHash, password)
	if err != nil {
		return nil, err
	}
	if !isValid {
		return nil, errors.New("invalid username or password")
	}
	account.PasswordHash = "" // Clear password hash before returning
	return account, nil
}

func ComparePasswords(hashedPassword, password string) (bool, error) {
	hashedPasswordBytes := []byte(hashedPassword)
	passwordBytes := []byte(password)
	err := bcrypt.CompareHashAndPassword(hashedPasswordBytes, passwordBytes)
	if err != nil {
		return false, err
	}
	return true, nil
}
