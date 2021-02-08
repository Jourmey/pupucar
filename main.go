package main

import (
	"flag"
	"fmt"
	"lockstepuiclient/client"
	"lockstepuiclient/ui"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten"
)

var mid = flag.Int("mid", 1, "help message for mid")
var roomid = flag.Int("roomid", 1, "help message for roomid")
var ip = flag.String("ip", "192.168.16.152", "lockstep server ip")

var roomID uint64
var id uint64
var u sync.Once

func main() {
	flag.Parse()

	rand.Seed(time.Now().UnixNano())
	roomID = uint64(*roomid)
	id = uint64(*mid)
	log.Print("welcome to ui ,mid = ", id, ", roomid = ", roomID)

	client.Run(roomID, id, *ip)
	show()
	//channel := client.GetLoginChan()
	//for i := range channel {
	//	switch i {
	//	case pb.ID_MSG_Connect:
	//	case pb.ID_MSG_Ready:
	//	case pb.ID_MSG_Start:
	//	}
	//}
}

func show() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle(fmt.Sprintf("ui id:%d room:%d", id, roomID))

	status := ui.NewStatus(id, roomID)
	g := ui.NewPupuCar(status)
	client.RegisterReceiveAction(
		g.ReceiveData,
		g.JoinRoom,
	)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
