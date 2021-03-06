package network

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	v1 "github.com/SotirisAlfonsos/chaos-bot/proto/grpc/v1"

	"github.com/go-kit/kit/log/level"

	"github.com/vishvananda/netlink"

	"github.com/SotirisAlfonsos/chaos-master/chaoslogger"
	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"
)

var (
	logger = getLogger()
)

type dataItem struct {
	message        string
	netemAttrs     *v1.NetworkRequest
	expectedResult *expected
}

type expected struct {
	message       string
	qdiscListSize int
}

func Test_Network_Delay_Injection_Success_Start(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	dataItems := []dataItem{
		{
			message:    "Start network delay with latency and duplication details",
			netemAttrs: &v1.NetworkRequest{Device: "lo", Latency: 1000, DelayCorr: 10.5, Duplicate: 10.5, DuplicateCorr: 10.5},
			expectedResult: &expected{
				message:       "Bot " + getHostname(t) + " started network injection with details {Latency: \\d+, Limit: 1000, Loss: 0, Gap: 0, Duplicate: \\d+, Jitter: 0}",
				qdiscListSize: 1,
			},
		},
		{
			message:    "Start network delay with network loss",
			netemAttrs: &v1.NetworkRequest{Device: "lo", Loss: 24.5},
			expectedResult: &expected{
				message:       "Bot " + getHostname(t) + " started network injection with details {Latency: 0, Limit: 1000, Loss: \\d+, Gap: 0, Duplicate: 0, Jitter: 0}",
				qdiscListSize: 1,
			},
		},
		{
			message: "Start network delay with full netem details",
			netemAttrs: &v1.NetworkRequest{Device: "lo", Latency: 10, DelayCorr: 10.2, Limit: 10, Loss: 10.2, LossCorr: 10.2, Gap: 10,
				Duplicate: 10.2, DuplicateCorr: 10.2, Jitter: 10, ReorderProb: 10.2, ReorderCorr: 10.2, CorruptProb: 10.2, CorruptCorr: 10.2},
			expectedResult: &expected{
				message:       "Bot " + getHostname(t) + " started network injection with details {Latency: \\d+, Limit: 10, Loss: \\d+, Gap: 10, Duplicate: \\d+, Jitter: \\d+}",
				qdiscListSize: 1,
			},
		},
	}

	for _, dataItem := range dataItems {
		runStartTest(t, dataItem)
	}
}

func runStartTest(t *testing.T, dataItem dataItem) {
	t.Run(dataItem.message, func(t *testing.T) {
		network := New(logger)

		resp, err := network.Start(dataItem.netemAttrs)
		defer cleanUp(t, network, dataItem.netemAttrs.Device)

		assert.Nil(t, err)
		assert.Regexp(t, regexp.MustCompile(dataItem.expectedResult.message), resp)
		assert.Equal(t, dataItem.expectedResult.qdiscListSize, getQdiscListSize(t, dataItem.netemAttrs.Device))
	})
}

func Test_Network_Delay_Injection_Error_Non_Existing_Link_Start(t *testing.T) {
	netemAttrs := &v1.NetworkRequest{
		Device: "non existing device",
	}

	network := New(logger)
	resp, err := network.Start(netemAttrs)

	assert.NotNil(t, err)
	assert.Regexp(t, fmt.Sprintf("Could not get link %s", netemAttrs.Device), resp)
}

func Test_Network_Delay_Injection_Error_Can_Not_Add_Two_Netem_For_A_Device(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	netemAttrs := &v1.NetworkRequest{
		Device:  "lo",
		Latency: 10,
	}

	network := New(logger)
	defer cleanUp(t, network, netemAttrs.Device)

	_, err := network.Start(netemAttrs)
	assert.Nil(t, err)

	resp2, err2 := network.Start(netemAttrs)
	assert.NotNil(t, err2)
	assert.Regexp(t, fmt.Sprintf("Could not add new qdisc %s", netemAttrs.Device), resp2)
}

func Test_Network_Delay_Injection_Success_Stop(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	dataItems := []dataItem{
		{
			message:    "Stop network delay with latency and duplication details",
			netemAttrs: &v1.NetworkRequest{Device: "lo", Latency: 10},
			expectedResult: &expected{
				message: "Bot " + getHostname(t) + " stopped network injection",
			},
		},
	}

	for _, dataItem := range dataItems {
		runStopTest(t, dataItem)
	}
}

func runStopTest(t *testing.T, dataItem dataItem) {
	t.Run(dataItem.message, func(t *testing.T) {
		network := New(logger)

		if _, err := network.Start(dataItem.netemAttrs); err != nil {
			t.Fatalf("failed to add qdisc netem for dev %s", dataItem.netemAttrs.Device)
		}

		resp, err := network.Stop(dataItem.netemAttrs.Device)

		assert.Nil(t, err)
		assert.Regexp(t, dataItem.expectedResult.message, resp)
	})
}

func Test_Network_Delay_Injection_Error_Non_Existing_Link_Stop(t *testing.T) {
	netemAttrs := &v1.NetworkRequest{
		Device: "non existing device",
	}

	network := New(logger)
	resp, err := network.Stop(netemAttrs.Device)

	assert.NotNil(t, err)
	assert.Regexp(t, fmt.Sprintf("Could not get link %s", netemAttrs.Device), resp)
}

func Test_Network_Delay_Injection_Error_Can_Not_Del_If_There_Is_No_Qdisc_Registered_For_Dev(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	netemAttrs := &v1.NetworkRequest{
		Device: "lo",
	}

	network := New(logger)

	resp, err := network.Stop(netemAttrs.Device)
	assert.NotNil(t, err)
	assert.Regexp(t, fmt.Sprintf("Could not delete qdisc %s", netemAttrs.Device), resp)
}

func getQdiscListSize(t *testing.T, device string) int {
	link, err := netlink.LinkByName(device)
	if err != nil {
		t.Fatalf("Could not get link %s", device)
	}
	qdiscs, err := netlink.QdiscList(link)
	if err != nil {
		t.Fatalf("Could not add new qdisc %s", device)
	}

	return len(qdiscs)
}
func cleanUp(t *testing.T, network *Network, device string) {
	_, err := network.Stop(device)
	if err != nil {
		t.Fatalf("Could not stop qdisc %s", device)
	}

	_ = level.Info(logger).Log("msg", "cleanup operation. Removed all network delays")
}

func getHostname(t *testing.T) string {
	hostname, err := os.Hostname()
	if err != nil {
		t.Fatalf("Can not get hostname")
	}
	return hostname
}

func getLogger() log.Logger {
	allowLevel := &chaoslogger.AllowedLevel{}
	if err := allowLevel.Set("debug"); err != nil {
		fmt.Printf("%v", err)
	}

	return chaoslogger.New(allowLevel)
}
