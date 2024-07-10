package user

import (
	"testing"
)

func TestGetAvatrURL(t *testing.T) {
	maps, err := getAvatrURL()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(maps)
}
