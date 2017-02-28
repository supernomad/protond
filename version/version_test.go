package version

import (
	"strconv"
	"strings"
	"testing"
)

func TestVersionIsSymantic(t *testing.T) {
	pieceies := strings.Split(GetVersion(), ".")

	if len(pieceies) != 3 {
		t.Fatal("version is not 3 segments")
	}

	i, err := strconv.Atoi(pieceies[0])
	if err != nil || i < 0 {
		t.Fatal("first segment is not a number, or segment is less than 0.")
	}

	i, err = strconv.Atoi(pieceies[1])
	if err != nil || i < 0 {
		t.Fatal("second segment is not a number, or segment is less than 0.")
	}

	i, err = strconv.Atoi(pieceies[2])
	if err != nil || i < 0 {
		t.Fatal("third segment is not a number, or segment is less than 0.")
	}
}
