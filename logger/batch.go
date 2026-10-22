package logger

import (
	"sync"
	"time"
)

type batchWriter struct {
	buf    []byte
	config OutputConfig
	mu     sync.Mutex
}

func newBatchWriter(oc OutputConfig) batchWriter {
	return batchWriter{
		buf:    make([]byte, 0, oc.BatchingFlushMaxBuffer),
		config: oc,
	}
}

func (bw *batchWriter) write(b []byte) {
	bw.mu.Lock()
	defer bw.mu.Unlock()
	if len(b)+len(bw.buf) >= bw.config.BatchingFlushMaxBuffer { // buffer limit is going to be hit with this append, then we flush first
		bw.flush()
	}
	bw.buf = append(bw.buf, b...)
}

// calls flush once per max flush time to flush output
func (bw *batchWriter) flusher() {
	t := time.NewTicker(bw.config.BatchingFlushMaxTime)
	for range t.C {
		bw.mu.Lock()
		bw.flush()
		bw.mu.Unlock()
	}
}

// non-mutexed write to output
func (bw *batchWriter) flush() {
	bw.config.Output.Write(bw.buf)
	bw.buf = bw.buf[:0] // resetting length
}
