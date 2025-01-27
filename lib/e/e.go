package e

import "fmt"

func Wrap(msg string, err error) error {
	if err != nil {
		return fmt.Errorf("%s: %w", msg, err) //%w — специальный формат для оборачивания ошибок
	}
	return nil
}

func WrapIfErr(msg string, err error) error {
	if err == nil {
		return nil
	}
	return Wrap(msg, err)
}
