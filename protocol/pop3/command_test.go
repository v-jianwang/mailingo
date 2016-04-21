package pop3

import (
	"testing"
)

func TestCommand(t *testing.T) {
	statString := "STAT"
	listString := "LIST 2"
	topString := "TOP 3 10"

	stat := NewCommand(statString)
	if stat.Name != "STAT" {
		t.Errorf("stat.Name parsed error: %+v", stat)
	}
	if len(stat.Args) != 0 {
		t.Errorf("stat.Args parsed error: %+v", stat)	
	}

	list := NewCommand(listString)
	if list.Name != "LIST" {
		t.Errorf("list.Name parsed error: %v", list)
	}
	if len(list.Args) != 1 {
		t.Errorf("list.Args parsed error: %v", list)	
	}
	if list.Args[0] != "2" {
		t.Errorf("list.Args[0] parsed error: %v", list)
	}

	top := NewCommand(topString)
	if top.Name != "TOP" {
		t.Errorf("top.Name parsed error: %v", top)
	}
	if len(top.Args) != 1 {
		t.Errorf("top.Args parsed error: %v", top)	
	}
	if top.Args[0] != "3" {
		t.Errorf("top.Args[0] parsed error: %v", top)
	}
	if top.Args[1] != "10" {
		t.Errorf("top.Args[1] parsed error: %v", top)
	}
}