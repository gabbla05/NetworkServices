package main

import (
	"fmt"
	"net"
	"os"
	"syscall"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Użycie: sudo go run injector.go [IP_DOCELOWE]")
		return
	}
	dstIP := os.Args[1]

	// 1. Tworzymy gniazdo surowe (Raw Socket) dla protokołu ICMP
	// AF_INET = IPv4, SOCK_RAW = surowe bajty, IPPROTO_ICMP = protokół kontrolny sieci
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_ICMP)
	if err != nil {
		fmt.Printf("Błąd: %v. Brak uprawnień do surowych gniazd (wymagane cap_net_raw).\n", err)
		return
	}
	defer syscall.Close(fd)

	// 2. Ręczne budowanie pakietu ICMP (Echo Request)
	// Typ: 8 (Echo), Kod: 0, Checksum: 0, ID: 1, Sequence: 1
	packet := []byte{
		8, 0, // Typ i Kod (ICMP Echo Request)
		0, 0, // Checksum (miejsce na sumę kontrolną)
		0, 1, // Identifier
		0, 1, // Sequence number
		'G', 'A', 'B', 'R', 'Y', // Dane (Twój podpis!)
	}

	// 3. Obliczanie sumy kontrolnej (wymagane w warstwie 3, żeby router nie odrzucił pakietu)
	cs := checksum(packet)
	packet[2] = byte(cs >> 8)
	packet[3] = byte(cs & 0xff)

	// 4. Adresowanie docelowe
	addr := syscall.SockaddrInet4{Port: 0}
	copy(addr.Addr[:], net.ParseIP(dstIP).To4())

	// 5. WYSYŁKA (Wstrzyknięcie bajtów prosto do karty sieciowej)
	err = syscall.Sendto(fd, packet, 0, &addr)
	if err != nil {
		fmt.Printf("Błąd wysyłania: %v\n", err)
	} else {
		fmt.Printf("✅ Sukces! Wysłano surowy pakiet ICMP (warstwa 3) do %s\n", dstIP)
		fmt.Printf("Wysłane dane (HEX): % x\n", packet)
	}
}

// Funkcja do obliczania sumy kontrolnej ICMP (standard RFC 1071)
func checksum(data []byte) uint16 {
	var sum uint32
	for i := 0; i < len(data)-1; i += 2 {
		sum += uint32(data[i])<<8 | uint32(data[i+1])
	}
	if len(data)%2 == 1 {
		sum += uint32(data[len(data)-1]) << 8
	}
	for sum > 0xffff {
		sum = (sum & 0xffff) + (sum >> 16)
	}
	return uint16(^sum)
}
