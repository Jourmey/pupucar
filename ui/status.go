package ui

type status struct {
	id     uint64
	roomid uint64
	level  int
}

func NewStatus(idd, roomidd uint64) *status {
	s := new(status)
	s.id = idd
	s.roomid = roomidd
	s.level = 0
	return s
}
