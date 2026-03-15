#!/bin/bash

# Skrypt do zarządzania ruchem sieciowym (iptables) na ocenę 5.0

PORT=80

case "$1" in
    block)
        echo "🚫 Blokuję ruch TCP/UDP na porcie $PORT..."
        # -A INPUT: Dodaj regułę do ruchu przychodzącego
        # -j DROP: Porzuć pakiet (klient dostanie Timeout)
        sudo iptables -A INPUT -p tcp --dport $PORT -j DROP
        sudo iptables -A INPUT -p udp --dport $PORT -j DROP
        echo "Zablokowano."
        ;;
    unblock)
        echo "✅ Odblokowuję ruch na porcie $PORT..."
        # -D INPUT: Usuń regułę
        sudo iptables -D INPUT -p tcp --dport $PORT -j DROP
        sudo iptables -D INPUT -p udp --dport $PORT -j DROP
        echo "Odblokowano."
        ;;
    status)
        echo "🔍 Aktualne reguły dla portu $PORT:"
        sudo iptables -L INPUT -n -v | grep $PORT || echo "Brak aktywnych blokad."
        ;;
    *)
        echo "Użycie: ./firewallSetup.sh [block|unblock|status]"
        ;;
esac