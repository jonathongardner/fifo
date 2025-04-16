package fifo

type ReadCloseReseter interface {
	Close() error
	Read(byte) (int, error)
	Reset() error
}
