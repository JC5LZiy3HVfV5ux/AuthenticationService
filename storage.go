package main

type Storage interface {
	Session() SessionManager
	Close()
}

type SessionManager interface {
	Insert(value *SessionDetails) error
	Get(guid string) (SessionDetails, error)
	Delete(guid string) error
}
