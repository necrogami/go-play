package fsm

import "fmt"
import "./server"


type Value interface {};

type Instance struct {
    server *Server;
    channel chan *Message;
    state func (Server, *Instance, *Message);
    states [](func (Server, *Instance, *Message));
}

type Message struct {
    What int;
    Data []Value;
}

type Server interface {
    Init(*Instance, *Message);
    HandleEvent(*Instance, *Message);
    HandleSyncEvent(*Instance, *Message);
}


func Start(srv Server, msg *Message) *Instance {
    inst := new(Instance);
    inst.server = &srv;
    inst.channel = make(chan *Message);

    go srv.Init(inst, msg);
    data := (<-inst.channel).Data;

    inst.states = make([](func (Server, *Instance, *Message)), len(data));
    for i, state := range data {
        inst.states[i] = state.(func (Server, *Instance, *Message));

        // Just set the first one as the initial state.
        if i == 0 {
            inst.state = inst.states[i];
        }
    }

    return inst;
}

func M(what int, data ...) *Message {
    msg := Message(*server.M(what, data));
    return &msg;
}


func (inst *Instance) Respond(msg *Message) {
    inst.channel <- msg;
}

func (inst *Instance) SetState(state func (Server, *Instance, *Message)) {
    inst.state = state;
}

func (inst *Instance) SendEvent(msg *Message) {
    go inst.state(*inst.server, inst, msg);
}

func (inst *Instance) SendSyncEvent(msg *Message) Value {
    go inst.server.HandleSyncEvent(inst, msg);
    result := <-inst.channel;

    if result.What != 0 {
        fmt.Printf("ERROR: %s\n", result.Data[0]);
    }

    return result;
}

func (inst *Instance) SendAllEvent(msg *Message) {
    go inst.server.HandleEvent(inst, msg);
}