package db

import (
	"testing"
)

func TestInClause(t *testing.T) {
	if InClause(1, 1) != "($1)" {
		t.Errorf("Failed. result wrong")
	}
	if InClause(1, 3) != "($3)" {
		t.Errorf("Failed. result wrong")
	}
	if InClause(2, 1) != "($1, $2)" {
		t.Errorf("Failed. result wrong")
	}
	if InClause(2, 3) != "($3, $4)" {
		t.Errorf("Failed. result wrong")
	}
	if InClause(5, 1) != "($1, $2, $3, $4, $5)" {
		t.Errorf("Failed. result wrong")
	}
	if InClause(5, 4) != "($4, $5, $6, $7, $8)" {
		t.Errorf("Failed. result wrong")
	}
}
