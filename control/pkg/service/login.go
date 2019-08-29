package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

type LoginObject struct {
	Account string `bson:"account,omitempty"`
	Password string `bson:"password,omitempty"`
}

func (l *LoginObject) Auth(id string, pw string) (err error) {
	cmp := LoginObject{
		Account: l.Account,
	}
	cmp.EncodePassword(id, pw)
	if strings.Compare(l.Password, cmp.Password) == 0 {
		return
	}
	return errors.New("password is incorrect")
}

func (l *LoginObject) ValidateAccount(account string) (err error) {
	if m, _ := regexp.MatchString(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`, account); m == false {
		err = fmt.Errorf("'account' is not email: %v", account)
		return
	}
	return
}

func (l *LoginObject) ValidatePassword(pw string) (err error){
	if len(pw) < 6 {
		err = errors.New("'password' is at least greater than 6")
		return
	}
	var num, lower, upper, spec bool
	for _, r := range pw {
		switch {
		case unicode.IsDigit(r):
			num = true
		case unicode.IsUpper(r):
			upper = true
		case unicode.IsLower(r):
			lower = true
		case unicode.IsSymbol(r), unicode.IsPunct(r):
			spec = true
		}
	}
	if num && lower && upper && spec {
		return nil
	}

	err = errors.New("'password' must contain uppercase, lowercase letters, numbers, and special characters")

	return
}

func (l *LoginObject) EncodePassword(id string, pw string) {
	l.Password = fmt.Sprintf("%x", sha1.Sum([]byte(pw + "?" + l.Account + "@" + id)))
}