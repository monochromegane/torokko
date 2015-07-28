package cargo

type storager interface {
	isExist() bool
	save(string, string) error
}

func newStorage(typ string) storager {
	switch typ {
	case "file":
		return newFileStorage()
	default:
		return nil
	}
}
