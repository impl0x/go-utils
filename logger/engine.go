package logger

import (
	"bytes"
	"math"

	"github.com/impl0x/go-utils/logger/utils"
)

type engine struct {
	config      Config
	pool        _pool
	timeWriter  timeWriter
	batchWriter batchWriter
}

func newEngine(config Config) *engine {
	n := func() *bytes.Buffer {
		return &bytes.Buffer{}
	}
	e := engine{
		config: config,
		pool: _pool{
			256,
			math.MaxInt,
			utils.NewPool(n),
		},
		timeWriter:  newTimeWriter(config.Time),
		batchWriter: newBatchWriter(config.Output),
	}
	go e.batchWriter.flusher()
	return &e
}

func (e *engine) log(l level, m string, kv []string) {
	if len(kv)%2 != 0 {
		panic("log: invalid number of key value pairs passed!")
	}
	b := e.pool.Get()
	writeColorStart(b, e.config.Color.Time)
	e.timeWriter.write(b)

	b.WriteString(e.config.Log.TimeLevelSep)

	writeColorStart(b, e.config.Color.Levels.fromLevel(l))
	writeLevel(b, l)

	b.WriteString(e.config.Log.LevelMsgSep)

	writeColorStart(b, e.config.Color.MessageColor)
	b.WriteString(m)

	if kv != nil {
		b.WriteString(e.config.Log.MessageExtraSep)
	}
	for i, v := range kv {
		if i%2 == 0 {
			writeColorStart(b, e.config.Color.ExKey)
			b.WriteString(v)
			b.WriteString(e.config.Log.ExtraKeyValSep)
		} else {
			writeColorStart(b, e.config.Color.ExVal)
			b.WriteString(v)
		}
	}
	writeColorEnd(b)

	b.WriteByte('\n')

	e.batchWriter.write(b.Bytes())

	b.Reset() // reset length

	// put back
	e.pool.Put(b)
}
