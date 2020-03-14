//-----------------------------------------------------------------------------
/*

SoC Interrupts

*/
//-----------------------------------------------------------------------------

package soc

//-----------------------------------------------------------------------------

// Interrupt describes an SoC interrupt.
type Interrupt struct {
	Name  string // name
	IRQ   uint   // interrupt request number
	Descr string // description
}

//-----------------------------------------------------------------------------