package interfaces

type ClientAccrual struct {
	addr string
}

func NewClienClientAccrual(addr string) *ClientAccrual {
	return &ClientAccrual{addr: addr}
}
