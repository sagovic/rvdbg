//-----------------------------------------------------------------------------
/*

Decode bitfields within registers.

*/
//-----------------------------------------------------------------------------

package soc

//-----------------------------------------------------------------------------

import (
	"fmt"
	"strings"

	"github.com/deadsy/rvdbg/util"
)

//-----------------------------------------------------------------------------

// Enum provides descriptive strings for the enumeration of a bit field.
type Enum map[uint]string

// fmtFunc is formatting function for a uint.
type fmtFunc func(x uint) string

// Field is a bit field within a register value.
type Field struct {
	Name       string  // name
	Msb        uint    // most significant bit
	Lsb        uint    // least significant bit
	Descr      string  // description
	Fmt        fmtFunc // formatting function
	Enums      Enum    // enumeration values
	cacheValid bool    // is the cached field value valid?
	cacheVal   uint    // cached field value
}

// FieldSet is a set of fields.
type FieldSet []Field

//-----------------------------------------------------------------------------
// Sort fields by Msb.

func (a FieldSet) Len() int      { return len(a) }
func (a FieldSet) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func (a FieldSet) Less(i, j int) bool {
	// MSBs for registers may not be unique.
	// Tie break with the name to give a well-defined sort order.
	if a[i].Msb == a[j].Msb {
		return strings.Compare(a[i].Name, a[j].Name) < 0
	}
	return a[i].Msb > a[j].Msb
}

//-----------------------------------------------------------------------------

// Display returns display strings for a bit field.
func (f *Field) Display(val uint) []string {
	// get the field
	val = util.Bits(val, f.Msb, f.Lsb)
	// has the value changed?
	changed := ""
	if val != f.cacheVal && f.cacheValid {
		changed = " *"
	}
	f.cacheVal = val
	f.cacheValid = true

	// field name
	var nameStr string
	if f.Msb == f.Lsb {
		nameStr = fmt.Sprintf("  %s[%d]", f.Name, f.Lsb)
	} else {
		nameStr = fmt.Sprintf("  %s[%d:%d]", f.Name, f.Msb, f.Lsb)
	}
	// value name
	valName := ""
	if f.Fmt != nil {
		valName = f.Fmt(val)
	} else if f.Enums != nil {
		if s, ok := f.Enums[val]; ok {
			valName = s
		}
	}
	// value string
	var valStr string
	if val < 10 {
		valStr = fmt.Sprintf(": %d %s%s", val, valName, changed)
	} else {
		valStr = fmt.Sprintf(": 0x%x %s%s", val, valName, changed)
	}
	return []string{nameStr, valStr, "", f.Descr}
}

//-----------------------------------------------------------------------------

// DisplayH returns the horizontal display string for the bit fields of a uint value.
func DisplayH(fs []Field, val uint) string {
	s := []string{}
	for _, f := range fs {
		x := util.Bits(val, f.Msb, f.Lsb)
		s = append(s, fmt.Sprintf("%s %s", f.Name, f.Fmt(x)))
	}
	return strings.Join(s, " ")
}

//-----------------------------------------------------------------------------
// standard formatting functions

// FmtDec formats a uint as a decimal string.
func FmtDec(x uint) string {
	return fmt.Sprintf("%d", x)
}

// FmtHex formats a uint as a hexadecimal string.
func FmtHex(x uint) string {
	return fmt.Sprintf("%x", x)
}

// FmtHex8 formats a uint as a 2-nybble hexadecimal string.
func FmtHex8(x uint) string {
	return fmt.Sprintf("%02x", x)
}

// FmtHex16 formats a uint as a 4-nybble hexadecimal string.
func FmtHex16(x uint) string {
	return fmt.Sprintf("%04x", x)
}

//-----------------------------------------------------------------------------
