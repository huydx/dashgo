package main

import (
	"testing"
	"strings"
)

func TestSearch(t *testing.T) {
	ms := manufactures()
	m := ms.find("88:57:ee:17:53:4c")
	if m == nil {
		t.Fatal("88:57:ee:17:53:4c should be found")
	}
	if strings.Contains(m.name, "Buffalo.Inc") != true {
		t.Fatal("88:57:ee:17:53:4c should be Apple")
	}
}
