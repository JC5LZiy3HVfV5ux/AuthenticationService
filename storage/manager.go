package storage

type Manager interface {
	Insert(value *Session) error
	Get(guid string) (*Session, error)
	Delete(guid string) error
}