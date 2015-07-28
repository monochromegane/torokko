package cargo

func newFileStorage() storager {
	return &fileStorage{}
}

type fileStorage struct {
}

func (f fileStorage) isExist() bool {
	return false
}

func (f fileStorage) save(dir, file string) error {
	return nil
}
