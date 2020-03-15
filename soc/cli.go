//-----------------------------------------------------------------------------
/*

SoC Device CLI

*/
//-----------------------------------------------------------------------------

package soc

import (
	"fmt"

	cli "github.com/deadsy/go-cli"
	"github.com/deadsy/rvdbg/util"
)

//-----------------------------------------------------------------------------

// Driver is the SoC driver api.
type Driver interface {
	GetAddressSize() uint
}

// target provides a method for getting the SoC device and driver.
type target interface {
	GetSoC() (*Device, Driver)
}

//-----------------------------------------------------------------------------

func getAddressFormat(drv Driver) string {
	return fmt.Sprintf("%%0%dx", drv.GetAddressSize()>>2)
}

//-----------------------------------------------------------------------------

var CmdMap = cli.Leaf{
	Descr: "display memory map",
	F: func(c *cli.CLI, args []string) {
		dev, drv := c.User.(target).GetSoC()
		s := make([][]string, len(dev.Peripherals))
		for i, p := range dev.Peripherals {
			var region string
			if p.Size == 0 {
				fmtStr := fmt.Sprintf(": %s", getAddressFormat(drv))
				region = fmt.Sprintf(fmtStr, p.Addr)
			} else {
				fmtStr := fmt.Sprintf(": %s %s %%s", getAddressFormat(drv), getAddressFormat(drv))
				region = fmt.Sprintf(fmtStr, p.Addr, p.Addr+p.Size-1, util.MemSize(p.Size))
			}
			s[i] = []string{p.Name, region, p.Descr}
		}
		c.User.Put(fmt.Sprintf("%s\n", cli.TableString(s, []int{0, 0, 0}, 1)))
	},
}

//-----------------------------------------------------------------------------

var RegsHelp = []cli.Help{
	{"<peripheral> [register]", "peripheral (string) - peripheral name"},
	{"", "register (string) - register name (or *)"},
}

var CmdRegs = cli.Leaf{
	Descr: "display peripheral registers",
	F: func(c *cli.CLI, args []string) {

		err := cli.CheckArgc(args, []int{1, 2})
		if err != nil {
			c.User.Put(fmt.Sprintf("%s\n", err))
			return
		}

		dev, drv := c.User.(target).GetSoC()

		p := dev.GetPeripheral(args[0])
		if p == nil {
			c.User.Put(fmt.Sprintf("no peripheral named \"%s\" (run \"map\" for the names)\n", args[0]))
			return
		}

		if len(p.Registers) == 0 {
			c.User.Put(fmt.Sprintf("peripheral \"%s\" has no registers\n", args[0]))
			return
		}

		if len(args) == 1 {
			c.User.Put(fmt.Sprintf("%s\n", p.Display(drv, nil, false)))
			return
		}
		if args[1] == "*" {
			c.User.Put(fmt.Sprintf("%s\n", p.Display(drv, nil, true)))
			return
		}

		r := p.GetRegister(args[1])
		if r == nil {
			c.User.Put(fmt.Sprintf("no register \"%s\" (run \"regs %s\" for the names) ", args[1], args[0]))
			return
		}
		c.User.Put(fmt.Sprintf("%s\n", p.Display(drv, r, true)))

	},
}

//-----------------------------------------------------------------------------
