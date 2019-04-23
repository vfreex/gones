package nes

type ClockRate int64

const (
	Herz      ClockRate = 1
	Kilohertz           = 1000 * Herz
	Megahertz           = 1000 * Kilohertz
	Gigahertz           = 1000 * Megahertz
)

// The clock rate of components in the NES
// http://wiki.nesdev.com/w/index.php/Cycle_reference_chart#Clock_rates
const (
	MasterClockRate = 236250 * Kilohertz / 11
	CpuClockRate    = MasterClockRate / 12
	PpuClockRate    = MasterClockRate / 4
)
