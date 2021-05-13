package config

type ES struct {
	Hosts   []string
	Indexes map[string]string
}

type Conn struct {
	ReqTimeoutMs int64
}

type InsertReq struct {
	Index string
	ID    string
}

func GetConn() *Conn {
	return &Conn{ReqTimeoutMs: 3000}
}

func (c *Conn) Insert(r InsertReq) {

}
