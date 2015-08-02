package main

type storager interface {
	isExist() bool
	save(string) error
	get(string) (string, error)
}

func newStorage(params *params) storager {
	switch storage {
	case "filesystem":
		return newFileStorage(params)
	default:
		return nil
	}
}
