package common

//Target is the interface with the commands for system and docker
type Target interface {
	Start(name string) (string, error)
	Stop(name string) (string, error)
}
