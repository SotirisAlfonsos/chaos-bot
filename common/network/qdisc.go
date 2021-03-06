package network

import (
	v1 "github.com/SotirisAlfonsos/chaos-bot/proto/grpc/v1"
	"github.com/vishvananda/netlink"
)

func getNetemQdisc(link netlink.Link, netemAttrs *v1.NetworkRequest) *netlink.Netem {
	qdiscAttributes := netlink.QdiscAttrs{
		LinkIndex: link.Attrs().Index,
		Handle:    netlink.MakeHandle(1, 0),
		Parent:    netlink.HANDLE_ROOT,
	}

	if netemAttrs == nil {
		return &netlink.Netem{
			QdiscAttrs: qdiscAttributes,
		}
	}

	netemAttributes := netlink.NetemQdiscAttrs{
		Latency:       netemAttrs.Latency * 1000,
		DelayCorr:     netemAttrs.DelayCorr,
		Limit:         netemAttrs.Limit,
		Loss:          netemAttrs.Loss,
		LossCorr:      netemAttrs.LossCorr,
		Gap:           netemAttrs.Gap,
		Duplicate:     netemAttrs.Duplicate,
		DuplicateCorr: netemAttrs.DuplicateCorr,
		Jitter:        netemAttrs.Jitter * 1000,
		ReorderProb:   netemAttrs.ReorderProb,
		ReorderCorr:   netemAttrs.ReorderCorr,
		CorruptProb:   netemAttrs.CorruptProb,
		CorruptCorr:   netemAttrs.CorruptCorr,
	}

	return netlink.NewNetem(qdiscAttributes, netemAttributes)
}
