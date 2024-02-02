package object

import "testing"

func TestStringHashKey(t *testing.T) {
	hello1 := &String{Value: "hello world"}
	hello2 := &String{Value: "hello world"}
	name1 := &String{Value: "yassinebenaid"}
	name2 := &String{Value: "yassinebenaid"}
	 
	if hello1.HashKey() != hello2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}

	if name1.HashKey() != name2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}

	if hello1.HashKey() == name1.HashKey() {
		t.Errorf("strings with different content have the same hash keys")
	}

	if hello2.HashKey() == name2.HashKey() {
		t.Errorf("strings with different content have the same hash keys")
	}
}

func TestIntegerHashKey(t *testing.T) {
	ten1 := &Integer{Value: 10}
	ten2 := &Integer{Value: 10}
	five1 := &Integer{Value: 5}
	five2 := &Integer{Value: 5}

	if ten1.HashKey() != ten2.HashKey() {
		t.Errorf("integers with same content have different hash keys")
	}

	if five1.HashKey() != five1.HashKey() {
		t.Errorf("integers with same content have different hash keys")
	}

	if ten1.HashKey() == five1.HashKey() {
		t.Errorf("integers with different content have the same hash keys")
	}

	if ten2.HashKey() == five2.HashKey() {
		t.Errorf("integers with different content have the same hash keys")
	}
}

func TestBooleanHashKey(t *testing.T) {
	true1 := &Boolean{Value: true}
	true2 := &Boolean{Value: true}
	false1 := &Boolean{Value: false}
	false2 := &Boolean{Value: false}

	if true1.HashKey() != true2.HashKey() {
		t.Errorf("booleans with same content have different hash keys")
	}

	if false1.HashKey() != false2.HashKey() {
		t.Errorf("booleans with same content have different hash keys")
	}

	if true1.HashKey() == false1.HashKey() {
		t.Errorf("booleans with different content have the same hash keys")
	}

	if true2.HashKey() == false2.HashKey() {
		t.Errorf("booleans with different content have the same hash keys")
	}
}
