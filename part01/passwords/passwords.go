package passwords

import (
	"golang.org/x/crypto/bcrypt"
)

type PasswordError struct {
	Msg string
}

func (p *PasswordError) Error() string {
	return p.Msg
}

// From https://medium.com/@jcox250/password-hash-salt-using-golang-b041dc94cb72

type PasswordStore struct {
	passwords map[string]string
}

func NewPasswordStore() *PasswordStore {
	p := PasswordStore{
		passwords: make(map[string]string),
	}
	return &p
}

func (p *PasswordStore) Add(uuid string, pwd string) error {

	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		return err
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it

	p.passwords[uuid] = string(hash)
	return nil
}

func (p *PasswordStore) Get(uuid string) (string, error) {
	hash, found := p.passwords[uuid]
	if !found {
		return "", &PasswordError{Msg: "Password not in system"}
	}
	return hash, nil
}

func (p *PasswordStore) Verify(uuid string, plainPwd string) (bool, error) {
	hashedPwd, getErr := p.Get(uuid)

	if getErr != nil {
		return false, getErr
	}

	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPwd))
	if err != nil {
		return false, &PasswordError{Msg: "Hash comparison failed"}
	}

	return true, nil
}
