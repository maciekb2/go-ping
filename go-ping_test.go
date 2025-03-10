package main

import (
	"bytes"
	"testing"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func TestCreateICMPEchoMessage(t *testing.T) {
	id := 1234
	seq := 1
	payload := "PING"

	packet, err := createICMPEchoMessage(id, seq, payload)
	if err != nil {
		t.Fatalf("Nieoczekiwany błąd przy tworzeniu pakietu: %v", err)
	}

	msg, err := icmp.ParseMessage(protocolICMP, packet)
	if err != nil {
		t.Fatalf("Błąd przy dekodowaniu pakietu: %v", err)
	}

	if msg.Type != ipv4.ICMPTypeEcho {
		t.Errorf("Oczekiwano ICMP Echo, otrzymano %v", msg.Type)
	}

	echo, ok := msg.Body.(*icmp.Echo)
	if !ok {
		t.Fatalf("Oczekiwano struktury icmp.Echo w ciele wiadomości")
	}

	if echo.ID != id {
		t.Errorf("Oczekiwano ID %d, otrzymano %d", id, echo.ID)
	}
	if echo.Seq != seq {
		t.Errorf("Oczekiwano Seq %d, otrzymano %d", seq, echo.Seq)
	}
	if !bytes.Equal(echo.Data, []byte(payload)) {
		t.Errorf("Oczekiwano danych %v, otrzymano %v", payload, echo.Data)
	}
}
