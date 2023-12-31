/*
Construido como parte da disciplina: Sistemas Distribuidos - PUCRS - Escola Politecnica
Professor: Fernando Dotti  (https://fldotti.github.io/)
Modulo representando Berst Effort Broadcast tal como definido em:

	Introduction to Reliable and Secure Distributed Programming
	Christian Cachin, Rachid Gerraoui, Luis Rodrigues

* Semestre 2018/2 - Primeira versao.  Estudantes:  Andre Antonitsch e Rafael Copstein
Para uso vide ao final do arquivo, ou aplicacao chat.go que usa este
*/
//package CausalOrderBroadcast
package CAUSAL_DIFUSION

import (
	"fmt"

	BestEffortBroadcast "SD/BEB"
)

type COBeB_Req_Message struct {
	Addresses []string
	W         []int //rel. vet
	Message   string
}

type COBeB_Ind_Message struct {
	From    string
	W       []int
	Message string
}

type COBEB_Module struct {
	Ind       chan COBeB_Ind_Message
	Req       chan COBeB_Req_Message
	id        int
	V         []int
	lsn       int //sequence number
	pending   []COBeB_Ind_Message
	addresses []string
	BEB       *BestEffortBroadcast.BestEffortBroadcast_Module
	dbg       bool
}

func (module *COBEB_Module) outDbg(s string) {
	if module.dbg {
		fmt.Println(". . . . . . . . . [ BEB msg : " + s + " ]")
	}
}

func (module *COBEB_Module) Init(addresses []string, id int, count int) {
	module.InitD(addresses, id, true, count)
}

func (module *COBEB_Module) InitD(addresses []string, id int, _dbg bool, count int) {
	module.dbg = _dbg
	module.outDbg("Init COBEB!")
	module.BEB = BestEffortBroadcast.NewBEB(addresses[id], id, _dbg)
	module.id = id
	module.addresses = addresses
	module.lsn = 0
	module.V = make([]int, count)
	for i := 0; i < count; i++ {
		module.V[i] = 0
	}
	module.Start()
}

func (module *COBEB_Module) Start() {

	go func() {
		for {
			select {
			case y := <-module.Req:
				module.handleBroadcast(y)
			case y := <-module.BEB.Ind:
				module.handleIndication(BEB2COBEB(y))
			}
		}
	}()
}

func (module *COBEB_Module) handleBroadcast(message COBeB_Req_Message) {
	W := module.V
	W[module.id] = module.lsn
	module.lsn = module.lsn + 1
	message.W = W
	module.BEB.Req <- COBEBReq2BEBReq(message)

}

func (module *COBEB_Module) handleIndication(message COBeB_Ind_Message) {
	module.pending = append(module.pending, message)
	for mIndex, msg := range module.pending {
		fromId := 0
		for i := 0; i < len(module.addresses); i++ {

			if module.addresses[i] == msg.From {
				fromId = i
			}
		}

		for clockIndex, _ := range msg.W {
			if msg.W[clockIndex] <= module.V[clockIndex] {
				module.pending = delete_at_index(module.pending, mIndex)
				module.V[fromId] = module.V[fromId] + 1
				fmt.Println("going to deliver")
				module.Deliver(msg)
				break
			}
		}
	}
}

func (module *COBEB_Module) Deliver(message COBeB_Ind_Message) {

	// fmt.Println("Received '" + message.Message + "' from " + message.From)
	module.Ind <- message
	// fmt.Println("# End BEB Received")
}

func COBEBReq2BEBReq(message COBeB_Req_Message) BestEffortBroadcast.BestEffortBroadcast_Req_Message {

	return BestEffortBroadcast.BestEffortBroadcast_Req_Message{
		Addresses: message.Addresses,
		W:         message.W,
		Message:   message.Message}

}

func BEB2COBEB(message BestEffortBroadcast.BestEffortBroadcast_Ind_Message) COBeB_Ind_Message {

	return COBeB_Ind_Message{
		From:    message.From,
		W:       message.W,
		Message: message.Message}

}

func NewCObeb(_id int, _addresses []string, _count int, _dbg bool) *COBEB_Module {

	cobeb := &COBEB_Module{
		Req: make(chan COBeB_Req_Message, 1),
		Ind: make(chan COBeB_Ind_Message, 1),
		id:  _id,
		dbg: _dbg,
	}
	cobeb.outDbg(" Init Causal Order Broadcast")
	cobeb.InitD(_addresses, _id, _dbg, _count)
	return cobeb
}

func delete_at_index(slice []COBeB_Ind_Message, index int) []COBeB_Ind_Message {
	if len(slice) == 1 {
		return []COBeB_Ind_Message{}
	} else {
		fmt.Println(len(slice))
		return append(slice[:index], slice[index+1:]...)
	}
}

// func main() {

// 	if len(os.Args) < 2 {
// 		fmt.Println("Please specify at least one address:port!")
// 		return
// 	}

// 	id, _ := strconv.Atoi(os.Args[1])
// 	addresses := os.Args[2:]
// 	fmt.Println(addresses)

// 	mod := COBEB_Module{
// 		Req: make(chan COBeB_Req_Message),
// 		Ind: make(chan COBeB_Ind_Message)}
// 	mod.Init(addresses, id, len(addresses))

// 	msg := COBeB_Req_Message{
// 		Addresses: addresses,
// 		Message:   "BATATA!"}

// 	yy := make(chan string)
// 	input := bufio.NewScanner(os.Stdin) //esperar
// 	input.Scan()
// 	mod.Req <- msg
// 	<-yy
// }
