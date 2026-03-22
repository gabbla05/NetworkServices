#!/bin/bash

echo "🚀 Budowanie projektów..."
go build -o server_bin ./p1-tcp-udp/server/server.go
go build -o injector_bin ./p1-raw-sockets/injector/injector.go
go build -o sniffer_bin ./p1-raw-sockets/sniffer/sniffer.go

echo "🛡️ Nadawanie Linux Capabilities (wymaga sudo tylko do nadania uprawnień plikom)..."

# Pozwala serwerowi działać na niskich portach (np. 80) bez roota
sudo setcap 'cap_net_bind_service=+ep' ./server_bin

# Pozwala na używanie gniazd surowych bez roota
sudo setcap 'cap_net_raw=+ep' ./injector_bin
sudo setcap 'cap_net_raw=+ep' ./sniffer_bin

echo "✅ Gotowe! Teraz możesz uruchamiać programy bez sudo:"
echo "1. Serwer na porcie 80: ./server_bin tcp 80"
echo "2. Injector: ./injector_bin 127.0.0.1"
echo "3. Sniffer: ./sniffer_bin"