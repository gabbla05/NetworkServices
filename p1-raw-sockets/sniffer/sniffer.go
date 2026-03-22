package main

import (
	"encoding/binary"
	"fmt"
	"syscall"
)

func main() {
	fmt.Println("🚀 Sniffer uruchomiony. Filtr: Port 80 (TCP i UDP).")
	fmt.Println("Czekam na pakiety...")

	// Uruchamiamy nasłuchiwanie TCP i UDP w osobnych wątkach (goroutines)
	go startSniffer(syscall.IPPROTO_TCP, "TCP")
	go startSniffer(syscall.IPPROTO_UDP, "UDP")

	// Zapobiega zamknięciu programu
	select {}
}

func startSniffer(protocol int, protoName string) {
	// Tworzymy gniazdo dla konkretnego protokołu
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, protocol)
	if err != nil {
		fmt.Printf("❌ Błąd [%s]: %v. Brak uprawnień cap_net_raw.\n", protoName, err)
		return
	}
	defer syscall.Close(fd)

	buf := make([]byte, 1500)
	for {
		n, _, err := syscall.Recvfrom(fd, buf, 0)
		if err != nil {
			continue
		}

		// Nagłówek IP ma 20 bajtów.
		// Zarówno w TCP, jak i w UDP port docelowy znajduje się na bajtach 22 i 23 bufora.
		if n > 23 {
			destPort := binary.BigEndian.Uint16(buf[22:24])

			// FILTR: Port 80
			if destPort == 80 {
				fmt.Printf("\n🎯 [%s] WYŁAPANO PAKIET DO SERWERA (Port: %d)\n", protoName, destPort)
				fmt.Printf("Długość całkowita: %d bajtów\n", n)
				fmt.Printf("Nagłówek IP (HEX): % x\n", buf[:20])

				if protoName == "TCP" {
					fmt.Printf("Nagłówek TCP (HEX): % x\n", buf[20:40])
				} else {
					fmt.Printf("Nagłówek UDP (HEX): % x\n", buf[20:28]) // UDP ma tylko 8 bajtów nagłówka
				}
				fmt.Println("--------------------------------------------------")
			}
		}
	}
}
