package cargo

type storager interface {
	isExist() bool
	save(string) error
}

func newStorage(typ string, params *params) storager {
	switch typ {
	case "file":
		return newFileStorage(params)
	default:
		return nil
	}
}
