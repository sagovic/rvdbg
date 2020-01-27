//-----------------------------------------------------------------------------
/*

Segger J-Link Driver

This package implements J-Link JTAG/SWD drivers using the libjaylink library.

*/
//-----------------------------------------------------------------------------

package jlink

import (
	"fmt"
	"strings"

	"github.com/deadsy/libjaylink"
	"github.com/deadsy/rvdbg/util/log"
)

//-----------------------------------------------------------------------------

const colorGreen = "\033[0;32m"
const colorNone = "\033[0m"

func logCallback(domain, msg string, user interface{}) {
	s := []string{colorGreen, domain, msg, colorNone}
	log.Debug.Printf("%s\n", strings.Join(s, ""))
}

//-----------------------------------------------------------------------------

// Jlink stores the J-Link library context.
type Jlink struct {
	ctx *libjaylink.Context
	dev []libjaylink.Device
}

// Init initializes the J-Link library.
func Init() (*Jlink, error) {
	// initialise the library
	ctx, err := libjaylink.Init()
	if err != nil {
		return nil, err
	}
	// setup the logging callback
	err = ctx.LogSetCallback(logCallback, nil)
	if err != nil {
		ctx.Exit()
		return nil, err
	}
	err = ctx.LogSetLevel(libjaylink.LOG_LEVEL_DEBUG)
	if err != nil {
		ctx.Exit()
		return nil, err
	}
	// discover devices
	err = ctx.DiscoveryScan(libjaylink.HIF_USB)
	if err != nil {
		ctx.Exit()
		return nil, err
	}
	dev, err := ctx.GetDevices()
	if err != nil {
		ctx.Exit()
		return nil, err
	}
	// return the library context
	j := &Jlink{
		ctx: ctx,
		dev: dev,
	}
	return j, nil
}

// Shutdown closes the J-Link library.
func (j *Jlink) Shutdown() {
	j.ctx.FreeDevices(j.dev, true)
	j.ctx.Exit()
}

// NumDevices returns the number of devices discovered.
func (j *Jlink) NumDevices() int {
	return len(j.dev)
}

// DeviceByIndex returns a J-Link device by index number.
func (j *Jlink) DeviceByIndex(idx int) (*libjaylink.Device, error) {
	if idx < 0 || idx >= len(j.dev) {
		return nil, fmt.Errorf("device index %d out of range", idx)
	}
	return &j.dev[idx], nil
}

//-----------------------------------------------------------------------------