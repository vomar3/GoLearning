package main

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"slices"
	"strings"
	"time"
	"unicode"
)

type User struct {
	FirstName         string
	LastName          string
	BirthYear         int
	FavoriteLanguages []string
}

func (u User) SecretIdentity() string {
	name := []rune(u.FirstName)
	surname := []rune(u.LastName)

	random := rand.IntN(100) + 1

	return fmt.Sprintf("%c%c%d", name[0], surname[0], random)
}

func (u User) Age() int {
	return time.Now().Year() - u.BirthYear
}

func (u *User) AddFavoriteLanguage(language string) error {
	if language == "" {
		return errors.New("empty language name")
	}

	if slices.Contains(u.FavoriteLanguages, language) {
		return errors.New("duplicate")
	}

	u.FavoriteLanguages = append(u.FavoriteLanguages, language)

	return nil
}

func (u *User) RemoveFavoriteLanguage(language string) error {
	if slices.Contains(u.FavoriteLanguages, language) {
		id := slices.Index(u.FavoriteLanguages, language)
		u.FavoriteLanguages = slices.Delete(u.FavoriteLanguages, id, id+1)
	} else {
		return errors.New("not found")
	}

	return nil
}

func (u User) IsProgrammingLanguageFavorite(language string) bool {
	return slices.Contains(u.FavoriteLanguages, language)
}

func (u User) RandomFavoriteLanguage() (string, error) {
	if len(u.FavoriteLanguages) == 0 {
		return "", errors.New("no options")
	}

	random := rand.IntN(len(u.FavoriteLanguages))

	return u.FavoriteLanguages[random], nil
}

func (u User) GenerateProfile() string {
	languages := "[" + strings.Join(u.FavoriteLanguages, ", ") + "]"

	return fmt.Sprintf("Имя: %s.\nФамилия: %s.\nВозраст: %d.\nСписок любимых языков программирования: %s.", u.FirstName, u.LastName, u.Age(), languages)
}

func (u *User) UpdateName(firstName, lastName string) error {
	if firstName == "" || lastName == "" {
		return errors.New("empty data")
	}

	name := []rune(firstName)
	surname := []rune(lastName)

	if !unicode.IsUpper(name[0]) || !unicode.IsUpper(surname[0]) {
		return errors.New("invalid data")
	}

	u.FirstName = firstName
	u.LastName = lastName

	return nil
}
