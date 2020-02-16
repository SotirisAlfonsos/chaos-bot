package common

import (
	"chaos-slave/proto"
)

type Target interface {
	Start() (*proto.StatusResponse, error)
	Stop() (*proto.StatusResponse, error)
}
