package common

// Target is the interface with the commands for system and docker
type Target interface {
	Start(item string) (string, error)
	Stop(item string) (string, error)
}
