/*
	Construido como parte da disciplina: Sistemas Distribuidos - PUCRS - Escola Politecnica
	Professor: Fernando Dotti  (https://fldotti.github.io/)
	Modulo representando Berst Effort Broadcast tal como definido em:
	  Introduction to Reliable and Secure Distributed Programming
	  Christian Cachin, Rachid Gerraoui, Luis Rodrigues
	* Semestre 2018/2 - Primeira versao.  Estudantes:  Andre Antonitsch e Rafael Copstein
	Para uso vide ao final do arquivo, ou aplicacao chat.go que usa este

// go run CausalDifuion_beb.go 0 127.0.0.1:5001  127.0.0.1:6001  127.0.0.1:7001 -> para testar, 3 terminais, com o parametro 0 e o ultimo de cada end. diferente
*/
package BookBestEffortBroadcast

import (
	"fmt"

	PP2PLink "SD/PP2PLink"
)

type BookBestEffortBroadcast_Req_Message struct {
	Message string
}

type BookBestEffortBroadcast_Ind_Message struct {
	From    int
	Message string
}

type BookBestEffortBroadcast_Module struct {
	Ind       chan BookBestEffortBroadcast_Ind_Message
	Req       chan BookBestEffortBroadcast_Req_Message
	id        int
	addresses []string
	Pp2plink  *PP2PLink.PP2PLink
	dbg       bool
}

func (module *BookBestEffortBroadcast_Module) outDbg(s string) {
	if module.dbg {
		fmt.Println(". . . . . . . . . [ BEB msg : " + s + " ]")
	}
}

func (module *BookBestEffortBroadcast_Module) Init(id int, addresses []string) {
	module.InitD(id, addresses, true)
}

func (module *BookBestEffortBroadcast_Module) InitD(id int, addresses []string, _dbg bool) {
	module.dbg = _dbg
	module.outDbg("Init BEB!")
	module.Pp2plink = PP2PLink.NewPP2PLink(addresses[id], _dbg)
	module.Start()
}

func (module *BookBestEffortBroadcast_Module) Start() {

	go func() {
		for {
			select {
			case y := <-module.Req:
				module.Broadcast(y)
			case y := <-module.Pp2plink.Ind:
				module.Deliver(PP2PLink2BEB(y, module.addresses))
			}
		}
	}()
}

func (module *BookBestEffortBroadcast_Module) Broadcast(message BookBestEffortBroadcast_Req_Message) {

	// aqui acontece o envio um opara um, para cada processo destinatario
	// em caso de injecao de falha no originador, no meio de um broadcast
	// este loop deve ser interrompido, tendo a mensagem ido para alguns mas nao para todos processos

	for i := 0; i < len(module.addresses); i++ {
		if i != module.id {
			msg := BEB2PP2PLink(message, i, module.addresses)
			//msg.To = module.addresses[i]
			module.Pp2plink.Req <- msg
			module.outDbg("Sent to " + module.addresses[i])
		}
	}
}

func (module *BookBestEffortBroadcast_Module) Deliver(message BookBestEffortBroadcast_Ind_Message) {

	// fmt.Println("Received '" + message.Message + "' from " + message.From)
	module.Ind <- message
	// fmt.Println("# End BEB Received")
}

func BEB2PP2PLink(message BookBestEffortBroadcast_Req_Message, toId int, addresses []string) PP2PLink.PP2PLink_Req_Message {

	return PP2PLink.PP2PLink_Req_Message{
		To:      addresses[toId],
		Message: message.Message}

}

func PP2PLink2BEB(message PP2PLink.PP2PLink_Ind_Message, addresses []string) BookBestEffortBroadcast_Ind_Message {
	fromID := 0
	for i := 0; i < len(addresses); i++ {
		if addresses[i] == message.From {
			fromID = i
		}
	}

	return BookBestEffortBroadcast_Ind_Message{
		From:    fromID,
		Message: message.Message}
}

func NewBookBEB(_id int, _addresses []string, _dbg bool) *BookBestEffortBroadcast_Module {
	beb := &BookBestEffortBroadcast_Module{
		Req: make(chan BookBestEffortBroadcast_Req_Message, 1),
		Ind: make(chan BookBestEffortBroadcast_Ind_Message, 1),

		dbg:      _dbg,
		Pp2plink: nil}
	beb.outDbg(" Init BestEffortBroadcast!")
	beb.Init(_id, _addresses)
	return beb
}

/*
func main() {

	if (len(os.Args) < 2) {
		fmt.Println("Please specify at least one address:port!")
		return
	}

	addresses := os.Args[1:]
	fmt.Println(addresses)

	mod := BookBestEffortBroadcast_Module{
		Req: make(chan BookBestEffortBroadcast_Req_Message),
		Ind: make(chan BookBestEffortBroadcast_Ind_Message) }
	mod.Init(addresses[0])

	msg := BookBestEffortBroadcast_Req_Message{
		Addresses: addresses,
		Message: "BATATA!" }

	yy := make(chan string)
	mod.Req <- msg
	<- yy
}
*/
