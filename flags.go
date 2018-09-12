package multi

import "log"

const (
	// Mcolor enables colored output.
	Mcolor = log.LUTC << (iota + 1) // doesn't conflict with log package flags
	// Mpid enables logging the process's PID.
	Mpid
)
