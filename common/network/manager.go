package network

import (
	"fmt"
	"os"

	v1 "github.com/SotirisAlfonsos/chaos-bot/proto/grpc/v1"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/vishvananda/netlink"
)

// Network is the interface implementation that injects network failures
type Network struct {
	Logger log.Logger
}

func New(logger log.Logger) *Network {
	return &Network{
		Logger: logger,
	}
}

func (n *Network) Start(netemAttrs *v1.NetworkRequest) (string, error) {
	link, err := netlink.LinkByName(netemAttrs.Device)
	if err != nil {
		return fmt.Sprintf("Could not get link %s", netemAttrs.Device), err
	}

	qdiscNetem := getNetemQdisc(link, netemAttrs)

	if err := netlink.QdiscAdd(qdiscNetem); err != nil {
		return fmt.Sprintf("Could not add new qdisc %s", netemAttrs.Device), err
	}

	return constructStartMessage(n.Logger, qdiscNetem), nil
}

func (n *Network) Recover(device string) (string, error) {
	link, err := netlink.LinkByName(device)
	if err != nil {
		return fmt.Sprintf("Could not get link %s", device), err
	}

	qdiscNetem := getNetemQdisc(link, nil)

	if err := netlink.QdiscDel(qdiscNetem); err != nil {
		return fmt.Sprintf("Could not delete qdisc %s", device), err
	}

	return constructMessage(n.Logger, "recovered"), nil
}

func constructStartMessage(logger log.Logger, qdiscNetem fmt.Stringer) string {
	message := constructMessage(logger, "started")
	return fmt.Sprintf("%s with details %s", message, qdiscNetem.String())
}

func constructMessage(logger log.Logger, action string) string {
	hostname, err := os.Hostname()
	if err != nil {
		_ = level.Warn(logger).Log("msg", "Could not get hostname", "err", err)
		hostname = "Unknown"
	}
	return fmt.Sprintf("Bot %s %s network injection", hostname, action)
}
