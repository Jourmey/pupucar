package main

import (
	"fmt"
	"github.com/byebyebruce/lockstepserver/protocol"
	"github.com/cihub/seelog"
	"github.com/golang/protobuf/proto"
	"github.com/xtaci/kcp-go"
	"lockstepclient/pb"
	"net"
	"time"
)

var logConfig = `
<seelog type="asynctimer" asyncinterval="5000000" minlevel="trace" maxlevel="error">
    <outputs formatid="common">
        <console/>
    </outputs>
    <formats>
        <format id="common" format="%Date %Time [%LEV] [%File:%Line] [%Func] %Msg%n"/>
    </formats>
</seelog>
`

var kk net.Conn
var room *pb.S2C_JoinRoomMsg
var roomseatid int
var frameID uint32
var frameID2 uint32
var sid int32
var isStart bool

func main() {
	initLog()
	initKcp()
	initUI()

	go read()
	go tick()
	err := sendMSG_Connect()
	if err != nil {
		panic(err)
	}

	select {}
}

func initUI() {

}

func initKcp() {
	k, err := kcp.Dial("localhost:10086")
	if err != nil {
		panic(err)
	}
	kk = k
}

func tick() {
	go func() {
		t1 := time.NewTicker(1 * time.Second)

		for range t1.C {
			if isStart {
				heartbeat()
			}
		}
	}()

	t := time.NewTicker(1 * time.Second / 60)

	for range t.C {
		if isStart {
			game()
		}
	}

}

func initLog() {
	li, _ := seelog.LoggerFromConfigAsString(logConfig)
	_ = seelog.ReplaceLogger(li)
	seelog.Info("seelog start success.")
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
				_ = seelog.Error("msg.UnmarshalPB failed. error=", err)
			} else {
				handleS2C_ConnectMsg(rec)
			}
		case pb.ID_MSG_Heartbeat:
			break
		case pb.ID_MSG_JoinRoom:
			rec := &pb.S2C_JoinRoomMsg{}
			if err := pp.UnmarshalPB(rec); nil != err {
				_ = seelog.Error("msg.UnmarshalPB failed. error=", err)
			} else {
				handleS2C_JoinRoomMsg(rec)
			}
		case pb.ID_MSG_Ready:
			seelog.Info("game ready!")
		case pb.ID_MSG_Start:
			seelog.Info("game start!")
			isStart = true
		case pb.ID_MSG_Frame:
			rec := &pb.S2C_FrameMsg{}
			if err := pp.UnmarshalPB(rec); nil != err {
				_ = seelog.Error("msg.UnmarshalPB failed. error=", err)
			} else {
				handleS2C_FrameMsg(rec)
			}
		default:
			_ = seelog.Warn("case id failed. id =", pp.GetMessageID())
		}
	}
}

func heartbeat() {
	_ = sendMsg(pb.ID_MSG_Heartbeat, nil)
	seelog.Info("heartbeat success.")
}

func game() {
	p := new(pb.C2S_InputMsg)
	p.X = proto.Int32(0)
	p.Y = proto.Int32(0)
	p.FrameID = proto.Uint32(frameID2)
	p.Sid = proto.Int32(sid)
	//p.X = proto.Int32(*p.X + 1)
	//p.Y = proto.Int32(*p.Y + 1)
	frameID2++
	sid++
	//if windev.KeyDownUp(windev.VK_UP) == 1 {
	//	p.Sid = proto.Int32(1)
	//} else if windev.KeyDownUp(windev.VK_DOWN) == 1 {
	//	p.Sid = proto.Int32(2)
	//} else if windev.KeyDownUp(windev.VK_LEFT) == 1 {
	//	p.Sid = proto.Int32(3)
	//} else if windev.KeyDownUp(windev.VK_RIGHT) == 1 {
	//	p.Sid = proto.Int32(4)
	//}

	_ = sendMsg(pb.ID_MSG_Input, p)

	seelog.Infof("sendMsg game data. (%d)\t(%d)\t(%d ,%d)", *p.FrameID, *p.Sid, *p.X, *p.Y)
}

func sendMSG_Connect() error {
	c := &pb.C2S_ConnectMsg{
		PlayerID: proto.Uint64(1),
		BattleID: proto.Uint64(1),
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
		seelog.Info("handleS2C_ConnectMsg success.")
		_ = sendMSG_JoinRoom()
	} else {
		_ = seelog.Warn("handleS2C_ConnectMsg failed. default rec = ", rec)
	}
}

func handleS2C_JoinRoomMsg(rec *pb.S2C_JoinRoomMsg) {
	seelog.Info("handleS2C_ConnectMsg success.")
	room = rec
	roomseatid = int(*rec.Roomseatid)
	_ = sendMSG_Ready()
}

func handleS2C_FrameMsg(rec *pb.S2C_FrameMsg) {
	seelog.Info("handleS2C_ConnectMsg success. rec = ", rec)
	if rec == nil {
		return
	}
	handleFrames(rec.Frames)
}

func handleFrames(frames []*pb.FrameData) {
	if frames == nil {
		return
	}
	for i := 0; i < len(frames); i++ {
		frameID = *frames[i].FrameID
		handleInputData(frames[i].Input)
	}
}

func handleInputData(input []*pb.InputData) {
	if input == nil || roomseatid >= len(input) {
		return
	}

	_ = input[roomseatid]
	//seelog.Infof("(%d)\t(%d)\t(%d)\t(%d ,%d)", *myInput.Id, *myInput.Sid, *myInput.Roomseatid, *myInput.X, *myInput.Y)
}
