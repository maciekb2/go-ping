/*
 * Program: go-ping (odpowiednik pinga systemowego)
 * Data utworzenia: 2025-03-10
 * Autor: mbrzezny
 *
 * Opis:
 *  Program wysyła pakiety ICMP Echo Request do wskazanego hosta i oczekuje na odpowiedzi.
 *  Obsługuje opcjonalne flagi:
 *      -host: adres IP lub nazwa hosta do pingowania (wymagane)
 *      -count: liczba wysyłanych pakietów (domyślnie 4)
 *      -timeout: czas oczekiwania na odpowiedź dla jednego pakietu (np. 3s)
 *      -interval: odstęp czasu między kolejnymi pingami (np. 1s)
 *      -payload: dane wysyłane w pakiecie ICMP (domyślnie "PING")
 *
 * Użycie:
 *  go run go-ping.go -host=8.8.8.8 -count=4 -timeout=3s -interval=1s -payload="PING"
 *
 * Uwaga:
 *  Na systemach Linux wymagane są uprawnienia administratora (root) lub CAP_NET_RAW.
 */

package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

const protocolICMP = 1

// createICMPEchoMessage tworzy pakiet ICMP Echo Request na podstawie przekazanych parametrów.
func createICMPEchoMessage(id, seq int, payload string) ([]byte, error) {
	echo := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   id,
			Seq:  seq,
			Data: []byte(payload),
		},
	}
	return echo.Marshal(nil)
}

func main() {
	// Definicja flag
	host := flag.String("host", "", "Adres IP lub nazwa hosta do pingowania (wymagane)")
	count := flag.Int("count", 4, "Liczba wysyłanych pakietów")
	timeout := flag.Duration("timeout", 3*time.Second, "Czas oczekiwania na odpowiedź dla jednego pakietu (np. 3s)")
	interval := flag.Duration("interval", 1*time.Second, "Odstęp czasu między kolejnymi pingami (np. 1s)")
	payload := flag.String("payload", "PING", "Dane wysyłane w pakiecie ICMP")
	flag.Parse()

	if *host == "" {
		fmt.Println("Błąd: Nie podano hosta. Użyj flagi -host")
		flag.Usage()
		os.Exit(1)
	}

	// Rozwiązywanie adresu IP hosta
	addr, err := net.ResolveIPAddr("ip4", *host)
	if err != nil {
		log.Fatalf("Błąd rozwiązywania adresu hosta: %v", err)
	}

	// Otwarcie połączenia ICMP
	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		log.Fatalf("Błąd przy otwieraniu połączenia ICMP: %v", err)
	}
	defer conn.Close()

	id := os.Getpid() & 0xffff
	var sent, received int
	var totalRTT time.Duration

	fmt.Printf("Pingowanie %s [%s] z %d bajtami danych:\n", *host, addr.String(), len(*payload))

	// Wysyłanie pakietów w pętli
	for i := 0; i < *count; i++ {
		packet, err := createICMPEchoMessage(id, i+1, *payload)
		if err != nil {
			log.Fatalf("Błąd przy tworzeniu pakietu ICMP: %v", err)
		}

		startTime := time.Now()
		_, err = conn.WriteTo(packet, addr)
		if err != nil {
			log.Printf("Błąd wysyłania pakietu: %v", err)
			continue
		}
		sent++

		// Ustawienie limitu oczekiwania na odpowiedź
		if err := conn.SetReadDeadline(time.Now().Add(*timeout)); err != nil {
			log.Fatalf("Błąd ustawiania limitu odczytu: %v", err)
		}

		// Odczyt odpowiedzi
		reply := make([]byte, 1500)
		n, peer, err := conn.ReadFrom(reply)
		rtt := time.Since(startTime)
		if err != nil {
			log.Printf("Przekroczono czas oczekiwania (timeout=%v): %v", *timeout, err)
		} else {
			receivedMsg, err := icmp.ParseMessage(protocolICMP, reply[:n])
			if err != nil {
				log.Printf("Błąd przy dekodowaniu pakietu: %v", err)
			} else {
				switch receivedMsg.Type {
				case ipv4.ICMPTypeEchoReply:
					// Próba parsowania ciała odpowiedzi
					echoReply, ok := receivedMsg.Body.(*icmp.Echo)
					if ok {
						fmt.Printf("Odpowiedź z %s: seq=%d czas=%v, odpowiedź: %s\n", peer, i+1, rtt, string(echoReply.Data))
					} else {
						fmt.Printf("Odpowiedź z %s: seq=%d czas=%v\n", peer, i+1, rtt)
					}
					received++
					totalRTT += rtt
				default:
					fmt.Printf("Nietypowa odpowiedź od %s: %+v\n", peer, receivedMsg)
				}
			}
		}

		// Odczekanie interwału przed kolejnym pingiem
		if i < *count-1 {
			time.Sleep(*interval)
		}
	}

	// Podsumowanie wyników
	fmt.Printf("\n--- Podsumowanie ping %s ---\n", *host)
	lost := sent - received
	loss := float64(lost) / float64(sent) * 100
	fmt.Printf("Pakiety: wysłane = %d, odebrane = %d, stracone = %d (%.1f%% utraconych)\n", sent, received, lost, loss)
	if received > 0 {
		avgRTT := totalRTT / time.Duration(received)
		fmt.Printf("Średni RTT: %v\n", avgRTT)
	}
}
