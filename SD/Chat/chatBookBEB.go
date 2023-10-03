// Construido como parte da disciplina: Sistemas Distribuidos - PUCRS - Escola Politecnica
//  Professor: Fernando Dotti  (https://fldotti.github.io/)

/*
LANCAR N PROCESSOS EM SHELL's DIFERENTES, UMA PARA CADA PROCESSO, O SEU PROPRIO ENDERECO EE O PRIMEIRO DA LISTA
go run chatBEB.go 127.0.0.1:5001  127.0.0.1:6001  127.0.0.1:7001   // o processo na porta 5001
go run chatBEB.go 127.0.0.1:6001  127.0.0.1:5001  127.0.0.1:7001   // o processo na porta 6001
go run chatBEB.go 127.0.0.1:7001  127.0.0.1:6001  127.0.0.1:5001     // o processo na porta ...
*/

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	. "SD/BookBEB"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Please specify at least one address:port!")
		fmt.Println("go run chatBookBEB.go 0 127.0.0.1:5001  127.0.0.1:6001   127.0.0.1:7001")
		fmt.Println("go run chatBookBEB.go 1 127.0.0.1:5001  127.0.0.1:6001   127.0.0.1:7001")
		fmt.Println("go run chatBookBEB.go 2 127.0.0.1:5001  127.0.0.1:6001   127.0.0.1:7001")
		return
	}

	var registro []string
	addresses := os.Args[2:]
	myId, _ := strconv.Atoi(os.Args[1])
	fmt.Println(addresses[myId])

	beb := NewBookBEB(myId, addresses, false)

	//beb.Init(addresses[0])
	//beb.InitD(myId, addresses, true)

	// enviador de broadcasts
	go func() {

		scanner := bufio.NewScanner(os.Stdin)
		var msg string

		for {
			if scanner.Scan() {
				msg = scanner.Text()
				msg += "§" + addresses[myId]
			}
			req := BookBestEffortBroadcast_Req_Message{
				Message: msg}
			beb.Req <- req // ENVIA PARA TODOS PROCESSOS ENDERECADOS NO INICIO
		}
	}()

	// receptor de broadcasts
	go func() {
		for {
			in := <-beb.Ind // RECEBE MENSAGEM DE QUALQUER PROCESSO
			message := strings.Split(in.Message, "§")
			in.Message = message[0]
			registro = append(registro, in.Message)

			// Adress correction
			for i := 0; i < len(addresses); i++ {

				if addresses[i] == message[1] {
					in.From = i
				}
			}

			// imprime a mensagem recebida na tela
			fmt.Printf("               Message from %v: %v\n", addresses[in.From], in.Message)
		}
	}()

	blq := make(chan int)
	<-blq
}
