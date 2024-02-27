package vm

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ControllerAdmissionMsg interface {
	sdk.Msg
	CheckAdmissibility(sdk.Context, interface{}) error

	// GetInboundMsgCount returns the number of Swingset messages which will
	// be added to the inboundQueue.
	GetInboundMsgCount() int32

	// IsHighPriority returns whether the message should be considered for
	// high priority processing, including bypass of some inbound checks
	// and queueing on higher priority queues.
	IsHighPriority(sdk.Context, interface{}) (bool, error)
}

var wrappedEmptySDKContext = sdk.WrapSDKContext(
	sdk.Context{}.WithContext(context.Background()),
)
var controllerContext context.Context = wrappedEmptySDKContext

type PortHandler interface {
	Receive(context.Context, string) (string, error)
}

var portToHandler = make(map[int]PortHandler)
var portToName = make(map[int]string)
var nameToPort = make(map[string]int)
var lastPort = 0

type protectedPortHandler struct {
	inner PortHandler
}

func (h protectedPortHandler) Receive(ctx context.Context, str string) (ret string, err error) {
	defer func() {
		if r := recover(); r != nil {
			// Propagate just the string, not the error stack or we will invite
			// nondeterminism.
			err = fmt.Errorf("panic: %s", r)
		}
	}()
	ret, err = h.inner.Receive(ctx, str)
	return
}

func NewProtectedPortHandler(inner PortHandler) PortHandler {
	return protectedPortHandler{inner}
}

func SetControllerContext(ctx sdk.Context) func() {
	// We are only called by the controller, so we assume that it is billing its
	// own meter usage.
	controllerContext = sdk.WrapSDKContext(ctx.WithGasMeter(sdk.NewInfiniteGasMeter()))
	return func() {
		controllerContext = wrappedEmptySDKContext
	}
}

func GetPort(name string) int {
	return nameToPort[name]
}

func RegisterPortHandler(name string, portHandler PortHandler) int {
	lastPort++
	portToHandler[lastPort] = NewProtectedPortHandler(portHandler)
	portToName[lastPort] = name
	nameToPort[name] = lastPort
	return lastPort
}
func UnregisterPortHandler(portNum int) error {
	delete(portToHandler, portNum)
	name := portToName[portNum]
	delete(portToName, portNum)
	delete(nameToPort, name)
	return nil
}
