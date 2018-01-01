package db

import (
	"encoding/json"
	"testing"
	"time"
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

type Person struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Birthday Time   `json:"birthday"`
}

func TestTimeJson(t *testing.T) {
	now := Time(time.Now())
	t.Log(now)
	src := `{"id":5,"name":"xiaoming","birthday":"2016-06-30 16:09:51"}`
	p := new(Person)
	err := json.Unmarshal([]byte(src), p)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(p)
	t.Log(time.Time(p.Birthday))
	js, _ := json.Marshal(p)
	t.Log(string(js))
}
