//-----------------------------------------------------------------------------
/*

GigaDevice gd32vf103 Flash Driver

This code implements the flash.Driver interface.

*/
//-----------------------------------------------------------------------------

package gd32vf103

import (
	"errors"
	"fmt"
	"time"

	"github.com/deadsy/rvdbg/mem"
	"github.com/deadsy/rvdbg/soc"
	"github.com/deadsy/rvdbg/util"
)

//-----------------------------------------------------------------------------

type flashMeta struct {
	name string
}

func (m *flashMeta) String() string {
	return m.name
}

//-----------------------------------------------------------------------------

// flashSectors returns a set of flash sectors for the device.
func flashSectors(dev *soc.Device) []*mem.Region {
	r := []*mem.Region{}
	// main flash
	p, _ := dev.GetPeripheral("flash")
	sectorSize := uint(1 * util.KiB)
	i := 0
	for addr := uint(p.Addr); addr < p.Addr+p.Size; addr += sectorSize {
		r = append(r, mem.NewRegion(p.Name, addr, sectorSize, &flashMeta{fmt.Sprintf("page %d", i)}))
		i++
	}
	// boot
	p, _ = dev.GetPeripheral("boot")
	r = append(r, mem.NewRegion(p.Name, p.Addr, p.Size, &flashMeta{"boot loader area"}))
	// option
	p, _ = dev.GetPeripheral("option")
	r = append(r, mem.NewRegion(p.Name, p.Addr, p.Size, &flashMeta{"option bytes"}))
	return r
}

//-----------------------------------------------------------------------------

// FlashDriver is a flash driver for the gd32vf103.
type FlashDriver struct {
	drv     soc.Driver
	dev     *soc.Device
	fmc     *soc.Peripheral
	sectors []*mem.Region
}

// NewFlashDriver returns a new gd32vf103 flash driver.
func NewFlashDriver(drv soc.Driver, dev *soc.Device) (*FlashDriver, error) {
	fmc, err := dev.GetPeripheral("FMC")
	if err != nil {
		return nil, err
	}
	return &FlashDriver{
		drv:     drv,
		dev:     dev,
		fmc:     fmc,
		sectors: flashSectors(dev),
	}, nil
}

// GetAddressSize returns the address size in bits.
func (drv *FlashDriver) GetAddressSize() uint {
	return 32
}

// GetDefaultRegion returns a default memory region.
func (drv *FlashDriver) GetDefaultRegion() *mem.Region {
	return mem.NewRegion("", 0, 1*util.KiB, nil)
}

// LookupSymbol returns an address and size for a symbol.
func (drv *FlashDriver) LookupSymbol(name string) *mem.Region {
	p, err := drv.dev.GetPeripheral(name)
	if err != nil {
		return nil
	}
	return mem.NewRegion(name, p.Addr, p.Size, nil)
}

// GetSectors returns the flash sector memory regions for the gd32vf103.
func (drv *FlashDriver) GetSectors() []*mem.Region {
	return drv.sectors
}

// Erase erases a flash sector.
func (drv *FlashDriver) Erase(r *mem.Region) error {
	time.Sleep(100 * time.Millisecond)
	return errors.New("TODO")
}

// EraseAll erases all of the device flash.
func (drv *FlashDriver) EraseAll() error {

	//# halt the cpu- don't try to run while we change flash
	//self.device.cpu.halt()

	//# make sure the flash is not busy
	//self.wait4complete()

	// unlock the flash
	err := drv.unlock()
	if err != nil {
		return err
	}

	//# set the mass erase bit
	err = drv.fmc.Set(drv.drv, "CTL0", ctlMER)
	if err != nil {
		return err
	}

	// set the start bit
	err = drv.fmc.Set(drv.drv, "CTL0", ctlSTART)
	if err != nil {
		return err
	}

	//# wait for completion
	//error = self.wait4complete()

	// clear the mass erase bit
	err = drv.fmc.Clr(drv.drv, "CTL0", ctlMER)
	if err != nil {
		return err
	}

	// lock the flash
	return drv.lock()
}

//-----------------------------------------------------------------------------
// private functions

const (
	ctlENDIE = (1 << 12) //                           End of operation interrupt enable bit
	ctlERRIE = (1 << 10) //                          Error interrupt enable bit
	ctlOBWEN = (1 << 9)  //                           Option byte erase/program enable bit
	ctlLK    = (1 << 7)  //                           FMC_CTL0 lock bit
	ctlSTART = (1 << 6)  //                           Send erase command to FMC bit
	ctlOBER  = (1 << 5)  //                           Option bytes erase command bit
	ctlOBPG  = (1 << 4)  //                           Option bytes program command bit
	ctlMER   = (1 << 2)  //                           Main flash mass erase for bank0 command bit
	ctlPER   = (1 << 1)  //                           Main flash page erase for bank0 command bit
	ctlPG    = (1 << 0)  // Main flash program for bank0 command bit
)

// unlock the flash
func (drv *FlashDriver) unlock() error {
	ctl, err := drv.fmc.Rd(drv.drv, "CTL0")
	if err != nil {
		return err
	}
	if ctl&ctlLK == 0 {
		// already unlocked
		return nil
	}
	// write the unlock sequence
	err = drv.fmc.Wr(drv.drv, "KEY0", 0x45670123)
	if err != nil {
		return err
	}
	err = drv.fmc.Wr(drv.drv, "KEY0", 0xCDEF89AB)
	if err != nil {
		return err
	}
	// clear any set CR bits
	err = drv.fmc.Wr(drv.drv, "CTL0", 0)
	if err != nil {
		return err
	}
	return nil
}

func (drv *FlashDriver) lock() error {
	return drv.fmc.Set(drv.drv, "CTL0", ctlLK)
}

/*

WS     : 40022000[31:0] = 0x00000030 wait state counter register
KEY0   : 40022004[31:0] = 0          Unlock key register 0
OBKEY  : 40022008[31:0] = 0          Option byte unlock key register
STAT0  : 4002200c[31:0] = 0          Status register 0
CTL0   : 40022010[31:0] = 0x00000080 Control register 0
ADDR0  : 40022014[31:0] = 0          Address register 0
OBSTAT : 4002201c[31:0] = 0x03fffffc Option byte status register
WP     : 40022020[31:0] = 0xffffffff Erase/Program Protection register
PID    : 40022100[31:0] = 0x4a425633 Product ID register

*/

//-----------------------------------------------------------------------------
