package files

//Этот пакет предоставляет реализацию хранилища на файловой системе.
//Данные сохраняются в директории, структура которой определяется пользователем.
//Сохраняет данные страницы (Page) на файловую систему.

import (
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"tg_bot/lib/e"
	"tg_bot/storage"
	"time"
)

type Storage struct {
	basePath string
}

const defaultPerm = 0774

//Определяет разрешения на создаваемые файлы и директории.
//0774: Чтение и запись для владельца и группы, чтение для остальных.

// Создаёт новый объект Storage с указанным базовым путём.
func New(basePath string) Storage {
	return Storage{basePath: basePath}
}

func (s Storage) Save(page *storage.Page) (err error) {
	//определяем способ обработки ошибок
	defer func() { err = e.WrapIfErr("can't save page", err) }()

	//определяем путь до дириктории, куда будет сохраняться файл
	fPath := filepath.Join(s.basePath, page.UserName)

	// Создаёт директорию, соответствующую имени пользователя
	if err := os.MkdirAll(fPath, defaultPerm); err != nil {
		return err
	}

	// Генерация имени файла
	fName, err := fileName(page)
	if err != nil {
		return err
	}

	//Добавляем к пути файла сгенерирванное имя
	fPath = filepath.Join(fPath, fName)

	//Создаем файл
	file, err := os.Create(fPath)
	if err != nil {
		return err
	}
	//Игнорируем ошибку о закрытии файла
	defer func() { _ = file.Close() }()

	//Записываем в файл страницу в нужном формате
	if err := gob.NewEncoder(file).Encode(page); err != nil {
		return err
	}

	return nil
}

func (s Storage) PickRandom(userName string) (page *storage.Page, err error) {
	defer func() { err = e.WrapIfErr("can't pick random page", err) }()

	//получаем путь до директории с файлами
	path := filepath.Join(s.basePath, userName)

	//1. проверить список папок внутри storage, если нет папки пользователя то не пытаемся проверить список файлов
	//просто сообщаем что ничего не сохранил

	//2.создать папку

	//получаем список файлов
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	//если файлов нет
	if len(files) == 0 {
		return nil, storage.ErrNoSavedPages
	}

	//генерация номера файла
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(files)) //верхняя граница

	file := files[n]

	//декодируем файл и возвращаем содержимое
	return s.decodePage(filepath.Join(path, file.Name()))
}

func (s Storage) Remove(p *storage.Page) error {
	fileName, err := fileName(p)
	if err != nil {
		return e.Wrap("can't remove file", err)
	}

	path := filepath.Join(s.basePath, p.UserName, fileName)

	if err := os.Remove(path); err != nil {
		msg := fmt.Sprintf("can't remove file %s", path)
		return e.Wrap(msg, err)
	}

	return nil
}

func (s Storage) IsExists(p *storage.Page) (bool, error) {
	fileName, err := fileName(p)
	if err != nil {
		return false, e.Wrap("can't check if file %s exists", err)
	}

	path := filepath.Join(s.basePath, p.UserName, fileName)

	switch _, err = os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		msg := fmt.Sprintf("can't check if file %s exists", path)

		return false, e.Wrap(msg, err)
	}

	return true, nil
}

func (s Storage) decodePage(filePath string) (*storage.Page, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, e.Wrap("can't decode page", err)
	}

	defer func() { _ = f.Close() }()

	//переменная, куда будет декодирован
	var p storage.Page

	if err := gob.NewDecoder(f).Decode(&p); err != nil {
		return nil, e.Wrap("can't decode page", err)
	}

	return &p, nil
}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}
