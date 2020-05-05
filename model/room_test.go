package model

import "testing"

func TestName(t *testing.T) {
	r := Room{}
	r.TablePieces.P1 = &TablePiecesOne{
		Pieces: Pieces{
			"1-1": 1,
			"1-2": 1,
		},
		Die: nil,
	}
	r.TablePieces.P2 = &TablePiecesOne{
		Pieces: Pieces{
			"1-0": 1,
			"1-3": 1,
		},
		Die: nil,
	}
	f, err := r.TablePieces.Move("p1", "1-1", "1-0")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%s %+v", f, r.TablePieces.P1)
}
