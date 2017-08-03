package steps

type repository interface {
	init() error
	clone() error
	commit() ([20]byte, error)
	push() error
	pull() error
}
