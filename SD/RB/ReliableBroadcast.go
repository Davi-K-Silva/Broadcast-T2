/*
Construido como parte da disciplina: Sistemas Distribuidos - PUCRS - Escola Politecnica
Professor: Fernando Dotti  (https://fldotti.github.io/)
Modulo representando Berst Effort Broadcast tal como definido em:

	Introduction to Reliable and Secure Distributed Programming
	Christian Cachin, Rachid Gerraoui, Luis Rodrigues

* Semestre 2018/2 - Primeira versao.  Estudantes:  Andre Antonitsch e Rafael Copstein
Para uso vide ao final do arquivo, ou aplicacao chat.go que usa este
*/
package ReliableBroadcast

import (
	"fmt"

	BookBestEffortBroadcast "SD/BookBEB"
)

type RB_Req_Message struct {
	Addresses []string
	Message   string
}

type RB_Ind_Message struct {
	From    int
	Message string
}

type ReliableBroadcast_Module struct {
	Ind       chan RB_Ind_Message
	Req       chan RB_Req_Message
	BookBEB   *BookBestEffortBroadcast.BookBestEffortBroadcast_Module
	id        int
	delivered []string

	dbg bool
}

func (module *ReliableBroadcast_Module) outDbg(s string) {
	if module.dbg {
		fmt.Println(". . . . . . . . . [ RB msg : " + s + " ]")
	}
}

func (module *ReliableBroadcast_Module) Init(addresses []string, id int) {
	module.InitD(addresses, id, true)
}

func (module *ReliableBroadcast_Module) InitD(addresses []string, id int, _dbg bool) {
	module.dbg = _dbg
	module.outDbg("Init RB!")
	module.id = id
	module.BookBEB = BookBestEffortBroadcast.NewBookBEB(id, addresses, _dbg)
	module.Start()
}

func (module *ReliableBroadcast_Module) Start() {

	go func() {
		for {
			select {
			case y := <-module.Req:
				module.Broadcast(y)
			case y := <-module.BookBEB.Ind:
				module.HandleDeliver(y)
			}
		}
	}()
}

func (module *ReliableBroadcast_Module) Broadcast(message RB_Req_Message) {
	module.BookBEB.Broadcast(RBReq2BEBReq(message))
}

func (module *ReliableBroadcast_Module) HandleDeliver(message BookBestEffortBroadcast.BookBestEffortBroadcast_Ind_Message) {
	if !stringInSlice(message.Message, module.delivered) {
		module.delivered = append(module.delivered, message.Message)
		module.Deliver(BEB2RBInd(message))
		module.BookBEB.Req <- BookBestEffortBroadcast.BookBestEffortBroadcast_Req_Message{
			Message: message.Message}
	}
}

func (module *ReliableBroadcast_Module) Deliver(message RB_Ind_Message) {

	// fmt.Println("Received '" + message.Message + "' from " + message.From)
	module.Ind <- message
	// fmt.Println("# End BookBEB Received")
}

func RBReq2BEBReq(message RB_Req_Message) BookBestEffortBroadcast.BookBestEffortBroadcast_Req_Message {

	return BookBestEffortBroadcast.BookBestEffortBroadcast_Req_Message{
		Message: message.Message}

}

func BEB2RBInd(message BookBestEffortBroadcast.BookBestEffortBroadcast_Ind_Message) RB_Ind_Message {

	return RB_Ind_Message{
		From:    message.From,
		Message: message.Message}

}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

/*
func main() {

	if (len(os.Args) < 2) {
		fmt.Println("Please specify at least one address:port!")
		return
	}

	addresses := os.Args[1:]
	fmt.Println(addresses)

	mod := ReliableBroadcast_Module{
		Req: make(chan BestEffortBroadcast_Req_Message),
		Ind: make(chan BestEffortBroadcast_Ind_Message) }
	mod.Init(addresses[0])

	msg := BestEffortBroadcast_Req_Message{
		Addresses: addresses,
		Message: "BATATA!" }

	yy := make(chan string)
	mod.Req <- msg
	<- yy
}
*/
