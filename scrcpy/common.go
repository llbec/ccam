package scrcpy

import (
	"fmt"
	"time"
)

//DebugLevel defines a varlue of debug level
type DebugLevel int

//DebugLevelMin - DebugLevelMax
const (
	DebugLevelMin DebugLevel = iota
	DebugLevelError
	DebugLevelWarn
	DebugLevelInfo
	DebugLevelDebug
	DebugLevelMax
)

//Debug return is debug mod
func (dl DebugLevel) Debug() bool {
	return dl >= DebugLevelDebug
}

//Info return is info mod
func (dl DebugLevel) Info() bool {
	return dl >= DebugLevelInfo
}

//Warn return is warn mod
func (dl DebugLevel) Warn() bool {
	return dl >= DebugLevelWarn
}

//Error return is error mod
func (dl DebugLevel) Error() bool {
	return dl >= DebugLevelError
}

//DebugLevelWrap int to level
func DebugLevelWrap(l int) DebugLevel {
	return DebugLevel(l % 6)
}

var debugOpt = DebugLevelMin

type size struct {
	width  uint16
	height uint16
}

func (s size) Center() Point {
	return Point{s.width >> 1, s.height >> 1}
}

func (s size) String() string {
	return fmt.Sprintf("size: (%d, %d)", s.width, s.height)
}

//Point struct x y
type Point struct {
	X uint16
	Y uint16
}

func (p Point) String() string {
	return fmt.Sprintf("Point: (%d, %d)", p.X, p.Y)
}

//PointMacro defines
type PointMacro struct {
	Point
	Interval time.Duration
}

//SPoint alias
type SPoint Point

//UserOperation interface
type UserOperation interface {
}
