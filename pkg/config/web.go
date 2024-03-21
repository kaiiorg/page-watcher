package config

var (
	DefaultHostPort = uint16(8080)
)

type Web struct {
	Port uint16 `hcl:"port"`
}

func defaultWeb() *Web {
	return &Web{
		Port: DefaultHostPort,
	}
}

func (w *Web) HostPort() uint16 {
	if w.Port == 0 {
		return DefaultHostPort
	}
	return w.Port
}
