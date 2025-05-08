package utils

import (
	"fmt"
	repo "forum/internal/repository"
	"regexp"
)

func InitRegex() {
	var err error
	repo.EmailExp, err = regexp.Compile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]+`)
	if err != nil {
		fmt.Println(err)
		return
	}

	repo.UsernameExp, err = regexp.Compile(`^[a-zA-Z]+[a-zA-Z0-9._]+`)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func ValidUsername(s string) bool {
	if len(s) < 8 || len(s) > 50 {
		return false
	}
	if !repo.UsernameExp.MatchString(s) {
		return false
	}
	return true
}

func ValidEmail(s string) bool {
	if len(s) < 8 || len(s) > 50 {
		return false
	}
	if !repo.EmailExp.MatchString(s) {
		return false
	}
	return true
}

func ValidPassword(s string) bool {
	if len(s) >= 8 && len(s) <= 50 {
		return true
	}
	return false
}
