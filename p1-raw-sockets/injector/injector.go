package main

import (
	"fmt"
	"net"
	"os"
	"syscall"
)

func main() {
	// Sprawdzamy argumenty: adres docelowy (z portem)
	if len(os.Args) < 2 {
		fmt.Println("Użycie: sudo go run injector.go [IP_DOCELOWE]")
		return
	}
	dstIP := os.Args[1]

	//GNIAZDO SUROWE DLA ICMP:
	// IPPROTO_ICMP mówi systemowi, że będziemy wysyłać pakiety kontrolne (typu Ping), a nie TCP czy UDP.
	// Wymaga to uprawnienia cap_net_raw
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
		0, 1, // Identifier: dowolny numer identyfikujący program, który wysyła pakiet
		0, 1, // Sequence number: kolejny numer pakietu
		'G', 'A', 'B', 'R', 'Y', // Dane
	}

	// 3. Obliczanie sumy kontrolnej (wymagane w warstwie 3, żeby router nie odrzucił pakietu)
	cs := checksum(packet)
	//wstawiamy wynik obliczen w drugim i trzecim bajcie pakietu (checksum zajmuje 2 bajty)
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
// sumuje 16bitowe kawalki danych, a potem bierze dopełnienie do jedynki (bitowe NOT)
func checksum(data []byte) uint16 {
	var sum uint32
	//przechodzimy po danych dwoma bajtami (16 bitów) i sumujemy je
	for i := 0; i < len(data)-1; i += 2 {
		sum += uint32(data[i])<<8 | uint32(data[i+1])
	}
	//jesli nieparzyscie dodajemy ostatni bajt jako 16-bitowy z zerem na końcu
	if len(data)%2 == 1 {
		sum += uint32(data[len(data)-1]) << 8
	}
	for sum > 0xffff {
		sum = (sum & 0xffff) + (sum >> 16)
	}
	//robimy negacjesumy, żeby uzyskać dopełnienie do jedynki (standard rfc)
	return uint16(^sum)
}
