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
		fmt.Println("Użycie: go run server.go [tcp|udp] [port]")
		return
	}

	protocol := strings.ToLower(os.Args[1])
	port := os.Args[2]

	if protocol == "tcp" {
		startTCPServer(port)
	} else if protocol == "udp" {
		startUDPServer(port)
	} else {
		fmt.Println("Nieznany protokół. Wybierz 'tcp' lub 'udp'.")
	}
}

func startTCPServer(port string) {
	// Słuchamy na wszystkich interfejsach. ":" oznacza Dual-Stack (IPv4 + IPv6)
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Printf("Błąd startu TCP: %v\n", err)
		return
	}
	defer ln.Close()
	fmt.Printf("Serwer TCP (IPv4/IPv6) nasłuchuje na porcie %s...\n", port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("Błąd akceptacji: %v\n", err)
			continue
		}
		go handleTCPConnection(conn)
	}
}

func handleTCPConnection(conn net.Conn) {
	defer conn.Close()
	// RemoteAddr() pokaże nam czy klient połączył się przez 127.0.0.1 czy [::1]
	fmt.Printf("[TCP] Nowe połączenie od: %s\n", conn.RemoteAddr().String())

	message, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Printf("[TCP] Odebrano: %s", message)
	conn.Write([]byte("Wiadomość odebrana pomyślnie!\n"))
}

func startUDPServer(port string) {
	addr, _ := net.ResolveUDPAddr("udp", ":"+port)
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Printf("Błąd startu UDP: %v\n", err)
		return
	}
	defer conn.Close()
	fmt.Printf("Serwer UDP nasłuchuje na porcie %s...\n", port)

	buf := make([]byte, 1024)
	for {
		n, remoteAddr, _ := conn.ReadFromUDP(buf)
		fmt.Printf("[UDP] Pakiet od %s: %s\n", remoteAddr, string(buf[:n]))
		conn.WriteToUDP([]byte("ACK: Paczka odebrana\n"), remoteAddr)
	}
}
