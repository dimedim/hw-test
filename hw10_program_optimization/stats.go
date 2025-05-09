package hw10programoptimization

import (
	"bufio"
	"io"
	"strings"
)

//go:generate easyjson -all stats.go
type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	scanner := bufio.NewScanner(r)
	suffix := "." + strings.ToLower(domain)
	resStat := make(DomainStat)

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var user User
		if err := user.UnmarshalJSON(line); err != nil {
			return nil, err
		}
		AtIdx := strings.LastIndex(user.Email, "@")
		if AtIdx < 0 {
			continue
		}
		domEmail := strings.ToLower(user.Email[AtIdx+1:])
		if strings.HasSuffix(domEmail, suffix) {
			resStat[domEmail]++
		}
	}
	return resStat, nil
}

// without easy json
// func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
// 	decoder := json.NewDecoder(r)
// 	suffix := "." + strings.ToLower(domain)
// 	resStat := make(DomainStat)
// 	for {
// 		var user User
// 		err := decoder.Decode(&user)
// 		if err == io.EOF {
// 			break
// 		} else if err != nil {
// 			return nil, err
// 		}
// 		AtIdx := strings.LastIndex(user.Email, "@")
// 		if AtIdx < 0 {
// 			continue
// 		}
// 		domEmail := strings.ToLower(user.Email[AtIdx+1:])
// 		if strings.HasSuffix(domEmail, suffix) {
// 			resStat[domEmail]++
// 		}
// 	}
// 	return resStat, nil
// }
