package shmipc

import (
	"encoding/binary"
	"errors"
	"testing"
)

func TestHeaderAndFallbackDataEvent(t *testing.T) {
	// Test header functions
	h := header(make([]byte, headerSize))
	h.encode(100, 2, typePolling)
	if h.Length() != 100 {
		t.Errorf("header Length() failed, expected %d, got %d", 100, h.Length())
	}
	if h.Magic() != magicNumber {
		t.Errorf("header Magic() failed, expected %d, got %d", magicNumber, h.Magic())
	}
	if h.Version() != 2 {
		t.Errorf("header Version() failed, expected %d, got %d", 1, h.Version())
	}
	if h.MsgType() != typePolling {
		t.Errorf("header MsgType() failed, expected %s, got %s", "Polling", h.MsgType().String())
	}

	// Test fallbackDataEvent functions
	f := fallbackDataEvent{}
	f.encode(200, 2, 12345, 67890)
	if binary.BigEndian.Uint32(f[0:4]) != 200 {
		t.Errorf("fallbackDataEvent encode() failed, expected %d, got %d", 200, binary.BigEndian.Uint32(f[0:4]))
	}
	if binary.BigEndian.Uint16(f[4:6]) != magicNumber {
		t.Errorf("fallbackDataEvent encode() failed, expected %d, got %d", magicNumber, binary.BigEndian.Uint16(f[4:6]))
	}
	if f[6] != 2 {
		t.Errorf("fallbackDataEvent encode() failed, expected %d, got %d", 1, f[6])
	}
	if f[7] != uint8(typeFallbackData) {
		t.Errorf("fallbackDataEvent encode() failed, expected %s, got %s", "FallbackData", eventType(f[7]).String())
	}
	if binary.BigEndian.Uint32(f[8:12]) != 12345 {
		t.Errorf("fallbackDataEvent encode() failed, expected %d, got %d", 12345, binary.BigEndian.Uint32(f[8:12]))
	}
	if binary.BigEndian.Uint32(f[12:16]) != 67890 {
		t.Errorf("fallbackDataEvent encode() failed, expected %d, got %d", 67890, binary.BigEndian.Uint32(f[12:16]))
	}
}

func TestCheckEventValid(t *testing.T) {
	validHeader := header{0x00, 0x00, 0x00, 0x10, 0x12, 0x34, 0x01, 0x02}
	invalidMagicHeader := header{0x00, 0x00, 0x00, 0x10, 0x56, 0x78, 0x01, 0x02}
	invalidVersionHeader := header{0x00, 0x00, 0x00, 0x10, 0x12, 0x34, 0x00, 0x02}
	invalidMsgTypeHeader := header{0x00, 0x00, 0x00, 0x10, 0x12, 0x34, 0x01, 0x0F}

	tests := []struct {
		name    string
		header  header
		wantErr error
	}{
		{"Valid header", validHeader, nil},
		{"Invalid magic header", invalidMagicHeader, errors.New("invalid magic or version")},
		{"Invalid version header", invalidVersionHeader, errors.New("invalid magic or version")},
		{"Invalid msg type header", invalidMsgTypeHeader, errors.New("invalid protocol header")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkEventValid(tt.header); err != tt.wantErr {
				t.Errorf("checkEventValid() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
