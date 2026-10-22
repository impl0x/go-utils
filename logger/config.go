package logger

import (
	"io"
	"os"
	"time"

	"github.com/impl0x/go-utils/logger/color"
)

type Option func(*Config)

// output config

func WithOutputWriter(w io.Writer) Option {
	return func(c *Config) {
		c.Output.Output = w
	}
}

func WithOutputBatchingOFF(c *Config) {
	c.Output.UseBatching = false
}

func WithOutputMaxFlushBufferSize(s uint) Option {
	return func(c *Config) {
		c.Output.BatchingFlushMaxBuffer = int(s)
	}
}

func WithOutputMaxFlushTime(td time.Duration) Option {
	return func(c *Config) {
		c.Output.BatchingFlushMaxTime = td
	}
}

// log config

func WithLogSeparators(lc LogConfig) Option {
	return func(c *Config) {
		c.Log = lc
	}
}

// time config

func WithTimeConfig(tc TimeConfig) Option {
	return func(c *Config) {
		c.Time = tc
	}
}

// color config

func WithColorConfig(cc ColorConfig) Option {
	return func(c *Config) {
		c.Color = cc
	}
}

func WithTimeColor(col color.RGB) Option {
	return func(c *Config) {
		c.Color.Time = col
	}
}

func WithMessageColor(col color.RGB) Option {
	return func(c *Config) {
		c.Color.MessageColor = col
	}
}

func WithExtraKeyColor(col color.RGB) Option {
	return func(c *Config) {
		c.Color.ExKey = col
	}
}

func WithExtraValueColor(col color.RGB) Option {
	return func(c *Config) {
		c.Color.ExVal = col
	}
}

func WithTraceColor(col color.RGB) Option {
	return func(c *Config) {
		c.Color.Levels.Trace = col
	}
}

func WithDebugColor(col color.RGB) Option {
	return func(c *Config) {
		c.Color.Levels.Debug = col
	}
}

func WithInfoColor(col color.RGB) Option {
	return func(c *Config) {
		c.Color.Levels.Info = col
	}
}

func WithWarnColor(col color.RGB) Option {
	return func(c *Config) {
		c.Color.Levels.Warn = col
	}
}

func WithErrorColor(col color.RGB) Option {
	return func(c *Config) {
		c.Color.Levels.Error = col
	}
}

func WithFatalColor(col color.RGB) Option {
	return func(c *Config) {
		c.Color.Levels.Fatal = col
	}
}

// other config

func WithPanicAtFatalOFF(c *Config) {
	c.Others.PanicAtFatal = false
}

type Config struct {
	Output OutputConfig
	Log    LogConfig
	Time   TimeConfig
	Color  ColorConfig
	Others OtherConfig
}

type OutputConfig struct {
	Output                 io.Writer
	UseBatching            bool          // uses batching to output logs else directly writes to output without batches on every call to logger
	BatchingFlushMaxTime   time.Duration // will hold the buffer for this long before flushing it to output
	BatchingFlushMaxBuffer int           // will hold the buffer till this length before flushing it to output
}

type LogConfig struct {
	TimeLevelSep    string // 09:24:10 PM" "INF Hello world extraKey=extraVal
	LevelMsgSep     string // 09:24:10 PM INF" "Hello world extraKey=extraVal
	MessageExtraSep string // 09:24:10 PM INF Hello world" "extraKey=extraVal
	ExtraKeyValSep  string // 09:24:10 PM INF Hello world extraKey" "extraVal
}

// config for the time part of the log,
// the time part is allocated 32 bytes only, so make sure the separators do not exceed too many bytes.
//
//   - date takes 8 or 6 bytes depending on if full year is displayed + date sep bytes
//   - datetime sep takes 1 according to default or whatever you choose
//   - time takes 6 + time sep bytes
//   - time format 12 hr takes extra 3 bytes due to space and AM PM, if space enabled else 2.
//     therefore it takes 17 bytes except all the separators, so make sure they take less than 15 bytes in *total*
type TimeConfig struct {
	TimeSep                    string // separator between time values, 12:29:23, here ":" is the separator
	DateSep                    string // separator between date values, 2026-02-02, here "-" is the separator
	DateTimeSep                string // separator between date and time, 2026-02-02 12:29:23, here " " is the separator
	ShowDate                   bool
	ShowFullYear               bool // 2026 / 26
	ShowTime                   bool
	TimeFormat                 TimeFormat
	AddSpaceFor12HrFormat      bool // 5 PM / 5PM
	AddZerosForSingleDigitDate bool // 26-04-01/26-4-1
	AddZerosForSingleDigitTime bool // 04:23 / 4:23
}

type TimeFormat uint8

const (
	Format24Hr TimeFormat = iota
	Format12Hr
)

type ColorConfig struct {
	Time         color.RGB
	Levels       LevelColors
	MessageColor color.RGB
	ExKey        color.RGB
	ExVal        color.RGB
}

type LevelColors struct {
	Trace color.RGB
	Debug color.RGB
	Info  color.RGB
	Warn  color.RGB
	Error color.RGB
	Fatal color.RGB
}

func (lc LevelColors) fromLevel(l level) color.RGB {
	switch l {
	case levelTrace:
		return lc.Trace
	case levelDebug:
		return lc.Debug
	case levelInfo:
		return lc.Info
	case levelWarn:
		return lc.Warn
	case levelError:
		return lc.Error
	case levelFatal:
		return lc.Fatal
	default:
		panic("invalid level passed")
	}
}

type OtherConfig struct {
	PanicAtFatal bool
}

var DefaultConfig = Config{
	Output: DefaultOutputConfig,
	Log:    DefaultLogConfig,
	Time:   DefaultTimeConfig,
	Color:  DefaultColorConfig,
	Others: DefaultOtherConfig,
}

var DefaultOutputConfig = OutputConfig{
	Output:                 os.Stdout,
	UseBatching:            true,
	BatchingFlushMaxTime:   250 * time.Millisecond, // will flush every 250 millisecond as default, decrease if want faster logging or increase if doesn't care about output timing, lowest can be a millisecond
	BatchingFlushMaxBuffer: 32768,                  // 32 kb buffer
}
var DefaultLogConfig = LogConfig{
	TimeLevelSep:    " ",
	LevelMsgSep:     " ",
	MessageExtraSep: " ",
	ExtraKeyValSep:  "=",
}
var DefaultTimeConfig = TimeConfig{
	TimeSep:                    ":",
	DateSep:                    "-",
	DateTimeSep:                " ",
	ShowDate:                   false,
	ShowFullYear:               true,
	ShowTime:                   true,
	TimeFormat:                 Format12Hr,
	AddSpaceFor12HrFormat:      true,
	AddZerosForSingleDigitDate: true,
	AddZerosForSingleDigitTime: true,
}
var DefaultLevelColors = LevelColors{
	Trace: color.NewRGB(128, 128, 128), // #808080
	Debug: color.NewRGB(0, 192, 255),   // #00c0ff
	Info:  color.NewRGB(0, 200, 0),     // #00c800
	Warn:  color.NewRGB(255, 165, 0),   // #ffa500
	Error: color.NewRGB(200, 0, 0),     // #c80000
	Fatal: color.NewRGB(255, 0, 255),   // #ff00ff
}
var DefaultColorConfig = ColorConfig{
	Time:         color.White[5],
	Levels:       DefaultLevelColors,
	MessageColor: color.White[1],
	ExKey:        color.NewRGB(10, 117, 119), // #0a7577
	ExVal:        color.White[3],
}
var DefaultOtherConfig = OtherConfig{
	PanicAtFatal: true,
}
