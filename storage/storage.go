package storage

//Этот пакет определяет интерфейс хранилища (Storage) и сущность страницы (Page).

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"read_adviser_tg_bot/lib/e"
)

type Storage interface {
	Save(p *Page) error
	PickRandom(userName string) (*Page, error)
	Remove(p *Page) error
	IsExists(p *Page) (bool, error)
}

var ErrNoSavedPages = errors.New("no saved pages")

type Page struct {
	URL      string
	UserName string
}

func (p Page) Hash() (string, error) {
	h := sha1.New()

	// Добавляет URL к хэшу
	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", e.Wrap("can't calculate hash", err)
	}

	// Добавляет имя пользователя к хэшу
	if _, err := io.WriteString(h, p.UserName); err != nil {
		return "", e.Wrap("can't calculate hash", err)
	}

	// Возвращает хэш в виде строки (шестнадцатеричный формат)
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
