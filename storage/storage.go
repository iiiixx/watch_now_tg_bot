package storage

type Storage interface {
	Save()
	PickRandom()
	Remove()
	IsExists()
}

type Page struct {
}
