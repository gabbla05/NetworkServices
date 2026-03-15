package main

import (
	"encoding/binary"
	"fmt"
	"syscall"
)

func main() {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_TCP)
	if err != nil {
		fmt.Printf("Błąd: %v. Brak uprawnień do surowych gniazd (wymagane cap_net_raw).\n", err)
		return
	}
	defer syscall.Close(fd)

	fmt.Println("Filtr ustawiony na port 80. Czekam na Twoje pakiety...")

	buf := make([]byte, 1500)
	for {
		n, _, err := syscall.Recvfrom(fd, buf, 0)
		if err != nil {
			continue
		}

		// Nagłówek IP ma zazwyczaj 20 bajtów.
		// Nagłówek TCP zaczyna się od 21. bajtu.
		// Port docelowy w TCP to bajty 22 i 23 (indeksy 22, 23 w buf).
		if n > 23 {
			destPort := binary.BigEndian.Uint16(buf[22:24])

			// FILTR: Pokaż tylko jeśli port docelowy to 8888
			if destPort == 80 {
				fmt.Printf("\n🎯 WYŁAPANO PAKIET DO SERWERA (Port: %d)\n", destPort)
				fmt.Printf("Długość całkowita: %d bajtów\n", n)
				fmt.Printf("Nagłówek IP (HEX): % x\n", buf[:20])
				fmt.Printf("Nagłówek TCP (HEX): % x\n", buf[20:40]) // Kolejne 20 bajtów to TCP
				fmt.Println("--------------------------------------------------")
			}
		}
	}
}
