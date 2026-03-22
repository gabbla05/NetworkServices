package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// Sprawdzamy argumenty: protokół i port
	if len(os.Args) < 3 {
		fmt.Println("Użycie: go run server.go [tcp|udp] [port]")
		return
	}

	protocol := strings.ToLower(os.Args[1])
	port := os.Args[2]

	//wybieramy odpowiednią funkcję serwera w zależności od protokołu
	if protocol == "tcp" {
		startTCPServer(port)
	} else if protocol == "udp" {
		startUDPServer(port)
	} else {
		fmt.Println("Nieznany protokół. Wybierz 'tcp' lub 'udp'.")
	}
}

// Funkcja serwera TCP, która obsługuje zarówno IPv4, jak i IPv6
func startTCPServer(port string) {
	// net.Listen: Otwieramy okienko i czekamy.
	// ":" przed portem to DUAL-STACK - serwer słucha jednocześnie na IPv4 i IPv6.
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Printf("Błąd startu TCP: %v\n", err)
		return
	}
	defer ln.Close()
	fmt.Printf("Serwer TCP (IPv4/IPv6) nasłuchuje na porcie %s...\n", port)

	for {
		// accept - czekamy na handshake od klienta i nawiazujemy połączenie. Zwraca nam obiekt conn, który reprezentuje to połączenie.
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("Błąd akceptacji: %v\n", err)
			continue
		}
		//obsluga klienta w osobnym watku (goroutine), żeby serwer mógł dalej nasłuchiwać i obsługiwać innych klientów
		go handleTCPConnection(conn)
	}
}

// obsluga konkretnego klienta TCP
func handleTCPConnection(conn net.Conn) {
	defer conn.Close() //po zakończeniu funkcji zamykamy połączenie
	// RemoteAddr() pokaże nam czy klient połączył się przez 127.0.0.1 czy [::1]
	fmt.Printf("[TCP] Nowe połączenie od: %s\n", conn.RemoteAddr().String())

	// Odczytujemy wiadomość od klienta (do znaku nowej linii)
	message, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Printf("[TCP] Odebrano: %s", message)
	conn.Write([]byte("Wiadomość odebrana pomyślnie!\n"))
}

// UDP bezpołączeniowy, więc nie ma handshake. Po prostu nasłuchujemy i odbieramy pakiety.
func startUDPServer(port string) {
	addr, _ := net.ResolveUDPAddr("udp", ":"+port) //dualstack
	//nie ma accept, po prostu nasluchujemy  na paczki przychodzace na dany port
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Printf("Błąd startu UDP: %v\n", err)
		return
	}
	defer conn.Close()
	fmt.Printf("Serwer UDP nasłuchuje na porcie %s...\n", port)

	buf := make([]byte, 1024) //miejsce na odebrane dane
	for {
		//ReadFromUDP - wyłapujemy paczkę, remoteAddr to adres nadawcy zeby wiedziec gdzie odpowiedziec
		n, remoteAddr, _ := conn.ReadFromUDP(buf)
		fmt.Printf("[UDP] Pakiet od %s: %s\n", remoteAddr, string(buf[:n]))
		//wysyłamy odpowiedź do nadawcy paczki
		conn.WriteToUDP([]byte("ACK: Paczka odebrana\n"), remoteAddr)
	}
}
