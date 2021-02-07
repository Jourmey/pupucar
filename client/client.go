package client

import (
	"fmt"
	"github.com/byebyebruce/lockstepserver/protocol"
	"github.com/golang/protobuf/proto"
	"github.com/xtaci/kcp-go"
	"lockstepuiclient/pb"
	"log"
	"net"
	"time"
)

var kk net.Conn

//var room *pb.S2C_JoinRoomMsg
var roomseatid int
var isStart bool
var StartChan = make(chan bool)

type InputDataEvent func(input []*pb.InputData)

var MHandler InputDataEvent

func Run(room, id uint64, serverAddr string) {
	initKcp(serverAddr)
	initUI()

	go read()
	go heart()
	err := sendMSG_Connect(room, id)
	if err != nil {
		panic(err)
	}
}

func SendAction(frameID uint32, sid int32) error {
	p := new(pb.C2S_InputMsg)
	p.FrameID = &frameID
	p.Sid = &sid
	return sendMsg(pb.ID_MSG_Input, p)
}

func initUI() {

}

func initKcp(add string) {
	k, err := kcp.Dial(fmt.Sprintf("%s:10086", add))
	if err != nil {
		panic(err)
	}
	kk = k
}

func heart() {
	go func() {
		t1 := time.NewTicker(1 * time.Second)
		for range t1.C {
			if isStart {
				heartbeat()
			}
		}
	}()
}

func read() {
	for {
		p := new(protocol.MsgProtocol)
		pack, err := p.ReadPacket(kk)
		if err != nil {
			panic(err)
		}
		pp := pack.(*protocol.Packet)

		switch pb.ID(pp.GetMessageID()) {
		case pb.ID_MSG_Connect:
			rec := &pb.S2C_ConnectMsg{}
			if err := pp.UnmarshalPB(rec); nil != err {
				log.Println("msg.UnmarshalPB failed. error=", err)
			} else {
				handleS2C_ConnectMsg(rec)
			}
		case pb.ID_MSG_Heartbeat:
			break
		case pb.ID_MSG_JoinRoom:
			rec := &pb.S2C_JoinRoomMsg{}
			if err := pp.UnmarshalPB(rec); nil != err {
				log.Println("msg.UnmarshalPB failed. error=", err)
			} else {
				handleS2C_JoinRoomMsg(rec)
			}
		case pb.ID_MSG_Ready:
			log.Print("game ready!")
		case pb.ID_MSG_Start:
			log.Print("game start!")
			isStart = true
			StartChan <- true
		case pb.ID_MSG_Frame:
			rec := &pb.S2C_FrameMsg{}
			if err := pp.UnmarshalPB(rec); nil != err {
				log.Println("msg.UnmarshalPB failed. error=", err)
			} else {
				handleS2C_FrameMsg(rec)
			}
		default:
			log.Print("case id failed. id =", pp.GetMessageID())
		}
	}
}

func heartbeat() {
	_ = sendMsg(pb.ID_MSG_Heartbeat, nil)
	log.Print("heartbeat success.")
}

func sendMSG_Connect(room, id uint64) error {
	c := &pb.C2S_ConnectMsg{
		PlayerID: proto.Uint64(id),
		BattleID: proto.Uint64(room),
		Token:    proto.String("token"),
	}
	return sendMsg(pb.ID_MSG_Connect, c)
}

func sendMSG_JoinRoom() error {
	return sendMsg(pb.ID_MSG_JoinRoom, nil)
}

func sendMSG_Ready() error {
	return sendMsg(pb.ID_MSG_Ready, nil)
}

func sendMsg(connect pb.ID, c interface{}) error {
	p := protocol.NewPacket(uint8(connect), c)
	if nil == p {
		return fmt.Errorf("p == nil")
	}

	_, err := kk.Write(p.Serialize())
	return err
}

func handleS2C_ConnectMsg(rec *pb.S2C_ConnectMsg) {
	if *rec.ErrorCode == pb.ERRORCODE_ERR_Ok {
		log.Print("handleS2C_ConnectMsg success.")
		_ = sendMSG_JoinRoom()
	} else {
		log.Print("handleS2C_ConnectMsg failed. default rec = ", rec)
	}
}

func handleS2C_JoinRoomMsg(rec *pb.S2C_JoinRoomMsg) {
	//room = rec
	roomseatid = int(*rec.Roomseatid)
	_ = sendMSG_Ready()
}

func handleS2C_FrameMsg(rec *pb.S2C_FrameMsg) {
	log.Print("handleS2C_ConnectMsg success. rec = ", rec)
	if rec == nil {
		return
	}
	handleFrames(rec.Frames)
}

func handleFrames(frames []*pb.FrameData) {
	if frames == nil {
		return
	}
	if MHandler != nil {
		for i := 0; i < len(frames); i++ {
			MHandler(frames[i].Input)
		}
	}
}
