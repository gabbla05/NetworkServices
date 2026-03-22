package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 3 {
		//sprawdzamy argumenty: protokół i adres docelowy (z portem)
		fmt.Println("Użycie: go run client.go [tcp|udp] [adres:port]")
		fmt.Println("Przykład IPv4: go run client.go tcp 127.0.0.1:8888")
		fmt.Println("Przykład IPv6: go run client.go tcp [::1]:8888")
		return
	}

	proto := strings.ToLower(os.Args[1]) //wybor protokołu
	addr := os.Args[2]                   //adres docelowy (z portem)

	if proto == "tcp" {
		//net.Dial: nawiązujemy połączenie z serwerem
		//jesli serwer slucha, to nastepuje handshake
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			fmt.Printf("Błąd połączenia TCP: %v\n", err)
			return
		}
		defer conn.Close() // po zakończeniu funkcji zamykamy połączenie

		// wyslanie wiadomosci do serwera i czekanie na odpowiedz
		fmt.Fprintf(conn, "Cześć serwerze przez TCP!\n")
		reply, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Printf("Serwer odpowiedział: %s", reply)

	} else if proto == "udp" {
		//net.Dial w UDP nie nawiazuje polaczenia tylko zapamietuje adres docelowy
		conn, err := net.Dial("udp", addr)
		if err != nil {
			fmt.Printf("Błąd UDP: %v\n", err)
			return
		}
		defer conn.Close()

		//rzucamy paczke z nadzieja ze serwer odbierze i odpowie
		fmt.Fprintf(conn, "Cześć serwerze przez UDP!\n")
		//miejsce na odp
		buf := make([]byte, 1024)
		n, _ := conn.Read(buf)
		fmt.Printf("Serwer odpowiedział: %s\n", string(buf[:n]))
	}
}
