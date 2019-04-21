package memory


type MemoryMapper interface {
	Map(addr Ptr) (mapped Ptr, ok bool)
}
