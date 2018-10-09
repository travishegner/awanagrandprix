package session

type Session interface {
	New() ([]byte, error)
	Load(sid []byte)
	Save() error
}
