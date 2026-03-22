package main

import (
	"encoding/binary"
	"fmt"
	"syscall"
)

func main() {
	fmt.Println("🚀 Sniffer uruchomiony. Filtr: Port 80 (TCP i UDP).")
	fmt.Println("Czekam na pakiety...")

	//uruchamiamy dwa sniffery naraz - jeden dla TCP, drugi dla UDP
	go startSniffer(syscall.IPPROTO_TCP, "TCP")
	go startSniffer(syscall.IPPROTO_UDP, "UDP")

	// Zapobiega zamknięciu programu
	select {}
}

func startSniffer(protocol int, protoName string) {
	// syscall.SOCK_RAW: pozwala na odbieranie surowych pakietów IP (z nagłówkiem)
	// wymaga uprawnienia cap_net_raw lub uruchomienia jako root
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, protocol)
	if err != nil {
		fmt.Printf("❌ Błąd [%s]: %v. Brak uprawnień cap_net_raw.\n", protoName, err)
		return
	}
	defer syscall.Close(fd)

	buf := make([]byte, 1500)
	for {
		//odbieramy surowy pakiet
		n, _, err := syscall.Recvfrom(fd, buf, 0)
		if err != nil {
			continue
		}

		// Nagłówek IP ma zawsze 20 bajtów. Zaraz po nim zaczyna się nagłówek TCP lub UDP.
		// Port docelowy to 2. i 3. bajt nagłówka warstwy transportowej.
		// Więc: 20 (IP) + 2 (offset portu) = 22. bajt w całym buforze.
		if n > 23 {
			//binary.BigEndian: zamienia 2 bajty na liczbe ludzka (port docelowy)
			destPort := binary.BigEndian.Uint16(buf[22:24])

			// FILTR interesują nas tylko pakiety, które są kierowane do portu 80
			if destPort == 80 {
				fmt.Printf("\n🎯 [%s] WYŁAPANO PAKIET DO SERWERA (Port: %d)\n", protoName, destPort)
				fmt.Printf("Długość całkowita: %d bajtów\n", n)
				fmt.Printf("Nagłówek IP (HEX): % x\n", buf[:20])

				if protoName == "TCP" {
					//min 20 bajtow naglowka tcp, ale moze byc wiecej (opcje), wiec pokazujemy 20-40 bajt
					fmt.Printf("Nagłówek TCP (HEX): % x\n", buf[20:40])
				} else {
					//UDP ma tylko 8 bajtów nagłówka, więc pokazujemy 20-28 bajt
					fmt.Printf("Nagłówek UDP (HEX): % x\n", buf[20:28]) // UDP ma tylko 8 bajtów nagłówka
				}
				fmt.Println("--------------------------------------------------")
			}
		}
	}
}
