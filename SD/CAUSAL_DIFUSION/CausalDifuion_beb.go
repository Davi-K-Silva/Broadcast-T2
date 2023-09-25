/*
Construido como parte da disciplina: Sistemas Distribuidos - PUCRS - Escola Politecnica
Professor: Fernando Dotti  (https://fldotti.github.io/)
Modulo representando Berst Effort Broadcast tal como definido em:

	Introduction to Reliable and Secure Distributed Programming
	Christian Cachin, Rachid Gerraoui, Luis Rodrigues

* Semestre 2018/2 - Primeira versao.  Estudantes:  Andre Antonitsch e Rafael Copstein
Para uso vide ao final do arquivo, ou aplicacao chat.go que usa este
*/
package CausalOrderBroadcast

import (
	"fmt"

	PP2PLink "SD/PP2PLink"
)

type BestEffortBroadcast_Req_Message struct {
	Addresses []string
	W         []int
	Message   string
}

type BestEffortBroadcast_Ind_Message struct {
	From    string
	W       []int
	Message string
}

type COBEB_Module struct {
	Ind      chan BestEffortBroadcast_Ind_Message
	Req      chan BestEffortBroadcast_Req_Message
	id       int
	V        []int
	lsn      int
	pending  []BestEffortBroadcast_Ind_Message
	Pp2plink *PP2PLink.PP2PLink
	dbg      bool
}

func (module *COBEB_Module) outDbg(s string) {
	if module.dbg {
		fmt.Println(". . . . . . . . . [ BEB msg : " + s + " ]")
	}
}

func (module *COBEB_Module) Init(address string, id int) {
	module.InitD(address, id, true)
}

func (module *COBEB_Module) InitD(address string, id int, _dbg bool) {
	module.dbg = _dbg
	module.outDbg("Init COBEB!")
	module.Pp2plink = PP2PLink.NewPP2PLink(address, _dbg)
	module.id = id
	module.lsn = 0
	for i := 0; i < len(module.V); i++ {
		module.V[i] = 0
	}
	module.Start()
}

func (module *COBEB_Module) Start() {

	go func() {
		for {
			select {
			case y := <-module.Req:
				module.Broadcast(y)
			case y := <-module.Pp2plink.Ind:
				module.Deliver(PP2PLink2BEB(y))
			}
		}
	}()
}

func handleBroadcast(module *COBEB_Module) {

}

func (module *COBEB_Module) Broadcast(message BestEffortBroadcast_Req_Message) {

	W := module.V
	W[module.id] = module.lsn
	module.lsn = module.lsn + 1
	message.W[0] = 1

	for i := 0; i < len(message.Addresses); i++ {
		msg := BEB2PP2PLink(message)
		msg.To = message.Addresses[i]
		module.Pp2plink.Req <- msg
		module.outDbg("Sent to " + message.Addresses[i])
	}
}

func (module *COBEB_Module) Deliver(message BestEffortBroadcast_Ind_Message) {

	// fmt.Println("Received '" + message.Message + "' from " + message.From)
	module.Ind <- message
	// fmt.Println("# End BEB Received")
}

func BEB2PP2PLink(message BestEffortBroadcast_Req_Message) PP2PLink.PP2PLink_Req_Message {

	return PP2PLink.PP2PLink_Req_Message{
		To:      message.Addresses[0],
		Message: message.Message}

}

func PP2PLink2BEB(message PP2PLink.PP2PLink_Ind_Message) BestEffortBroadcast_Ind_Message {

	return BestEffortBroadcast_Ind_Message{
		From:    message.From,
		Message: message.Message}
}

/*
func main() {

	if (len(os.Args) < 2) {
		fmt.Println("Please specify at least one address:port!")
		return
	}

	addresses := os.Args[1:]
	fmt.Println(addresses)

	mod := COBEB_Module{
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
