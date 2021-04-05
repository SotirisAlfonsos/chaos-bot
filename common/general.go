package common

// Target is the interface with the commands for system and docker
type Target interface {
	Recover(item string) (string, error)
	Kill(item string) (string, error)
}
