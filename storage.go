package goeinstein

type Storage interface {
	GetInt(name string, dflt int) int
	GetString(name string, dflt string) string
	SetInt(name string, value int)
	SetString(name string, value string)
	Flush()
	Close()
}

type StorageHolder struct {
	storage Storage
}

func (s *StorageHolder) GetStorage() Storage { return s.storage }

func NewStorageHolder() *StorageHolder {
	s := &StorageHolder{}
	s.storage = NewTableStorage()
	return s
}

func (s *StorageHolder) Close() {
	if s.storage != nil {
		s.storage.Close()
	}
}

var storageHolder = NewStorageHolder()

func GetStorage() Storage { return storageHolder.GetStorage() }
