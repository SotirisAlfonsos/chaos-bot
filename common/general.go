package common

type Target interface {
	Start(name string) (string, error)
	Stop(name string) (string, error)
}
