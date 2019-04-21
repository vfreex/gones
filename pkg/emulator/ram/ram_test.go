package ram

import (
	"github.com/vfreex/gones/pkg/emulator/memory"
	"testing"
)

func TestRamWriteRead(t *testing.T) {
	r := NewRAM(1024)
	ptr := memory.Ptr(0)
	e := byte(234)
	r.Poke(ptr, e)
	v := r.Peek(ptr)
	if v != e {
		t.Error("For", ptr, "expected", e, "got", v)
	}
}
