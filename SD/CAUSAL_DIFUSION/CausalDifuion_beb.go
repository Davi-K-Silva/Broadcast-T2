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
package main

import (
	"fmt"
	"strconv"
	"os"
	"bufio"

	PP2PLink "SD/PP2PLink"
)

type COBeB_Req_Message struct {
	Addresses []string
	W         []int//rel. vet
	Message   string
}

type COBeB_Ind_Message struct {
	From    string
	W       []int
	Message string
}

type COBEB_Module struct {
	Ind      chan COBeB_Ind_Message
	Req      chan COBeB_Req_Message
	id       int
	V        []int
	lsn      int //sequence number
	pending  []COBeB_Ind_Message
	Pp2plink *PP2PLink.PP2PLink
	dbg      bool
}

func (module *COBEB_Module) outDbg(s string) {
	if module.dbg {
		fmt.Println(". . . . . . . . . [ BEB msg : " + s + " ]")
	}
}

func (module *COBEB_Module) Init(address string, id int, count int) {
	module.InitD(address, id, true, count)
}

func (module *COBEB_Module) InitD(address string, id int, _dbg bool, count int) {
	module.dbg = _dbg
	module.outDbg("Init COBEB!")
	module.Pp2plink = PP2PLink.NewPP2PLink(address, _dbg)
	module.id = id
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
			case y := <-module.Ind:
				module.handleIndication(y)
			}
		}
	}()
}

func (module *COBEB_Module) handleBroadcast(message COBeB_Req_Message) {
	W := module.V
	W[module.id] = module.lsn
	module.lsn = module.lsn + 1
	message.W = W
	message.Message = strconv.Itoa(module.id)
	module.Broadcast(message)
}

func (module *COBEB_Module) handleIndication(message COBeB_Ind_Message) {
	module.pending = append(module.pending, message)
	for _, msg := range module.pending {
		fromId, _ := strconv.Atoi(msg.Message)
		for clockIndex, _ := range msg.W {
			if msg.W[clockIndex] <= module.V[clockIndex] {
				module.V[fromId] = module.V[fromId] + 1
				module.Deliver(msg)
			}
		}
	}
}

func (module *COBEB_Module) Broadcast(message COBeB_Req_Message) {

	for i := 0; i < len(message.Addresses); i++ {
		msg := COBEBReq2PP2PLinkReq(message)
		msg.To = message.Addresses[i]
		module.Pp2plink.Req <- msg
		module.outDbg("Sent to " + message.Addresses[i])
	}
}

func (module *COBEB_Module) Deliver(message COBeB_Ind_Message) {

	// fmt.Println("Received '" + message.Message + "' from " + message.From)
	module.Ind <- message
	// fmt.Println("# End BEB Received")
}

func COBEBReq2PP2PLinkReq(message COBeB_Req_Message) PP2PLink.PP2PLink_Req_Message {

	return PP2PLink.PP2PLink_Req_Message{
		To:      message.Addresses[0],
		Message: message.Message}

}

func main() {

	if (len(os.Args) < 2) {
		fmt.Println("Please specify at least one address:port!")
		return
	}

	id, _ := strconv.Atoi(os.Args[1])
	addresses := os.Args[2:]
	fmt.Println(addresses)

	mod := COBEB_Module{
		Req: make(chan COBeB_Req_Message),
		Ind: make(chan COBeB_Ind_Message) }
	mod.Init(addresses[id], id, len(addresses))

	msg := COBeB_Req_Message{
		Addresses: addresses,
		Message: "BATATA!" }

	yy := make(chan string)
	input:= bufio.NewScanner(os.Stdin)//esperar
	input.Scan()
	mod.Req <- msg
	<- yy
}

