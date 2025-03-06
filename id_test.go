package souvenir

import "testing"

type TestType struct{}

func (t TestType) Prefix() string { return "u" }

func TestZero(t *testing.T) {
	id := ZeroID[TestType]()
	e := "u_00000000000000000000000000"

	if id.String() != e {
		t.Errorf(`Expected ZeroID().String() = %s, got %s`, e, id.String())
	}
}

func TestBytes(t *testing.T) {
	id, _ := ParseID[TestType]("u_3456789abc0000123456789abc")
	bytes := [16]byte{0x64, 0x29, 0x8e, 0x84, 0xa9, 0x6c, 0x00, 0x00, 0x00, 0x88, 0x64, 0x29, 0x8e, 0x84,
		0xa9, 0x6c}
	if id.data != bytes {
		t.Errorf(`Expected %v but got %v`, bytes, id.data)
	}

	id2, _ := ParseID[TestType]("u_7zzzzzzzzzzzzzzzzzzzzzzzzz")
	bytes2 := [16]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff}
	if id2.data != bytes2 {
		t.Errorf(`Expected %v but got %v`, bytes2, id2.data)
	}

}

func TestRoundTrip(t *testing.T) {
	for range 1_000_000 {
		id := RandomID[TestType]()
		str := id.String()
		id2, err := ParseID[TestType](str)

		if err != nil {
			t.Errorf(`ParseID error: %q`, err)
		}

		if id.data != id2.data {
			t.Errorf(`Round trip mismatch: %q %q`, id.data, id2.data)
		}

		uuid := id.UUID()
		id3 := ParseUUID[TestType](uuid)

		if id3.data != id.data {
			t.Errorf(`Round trip uuid mismatch: %q %q`, id.data, id2.data)
		}
	}
}

