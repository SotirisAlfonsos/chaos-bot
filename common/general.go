package common

//Target is the interface with the commands for system and docker
type Target interface {
	Start() (string, error)
	Stop() (string, error)
}
