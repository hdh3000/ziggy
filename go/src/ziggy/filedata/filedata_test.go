package filedata

import (
	"log"
	"testing"
)

type Dummy struct {
	Name       string `fd:"key"`
	OtherThing int
}

func TestMgr_Put(t *testing.T) {
	mgr := NewMgr("./tmp.json")

	if err := mgr.Put(&Dummy{Name: "foo", OtherThing: 1}); err != nil {
		log.Fatal(err)
	}
	if err := mgr.Put(&Dummy{Name: "bar", OtherThing: 2}); err != nil {
		log.Fatal(err)
	}
	if err := mgr.Put(&Dummy{Name: "foo", OtherThing: 4}); err != nil {
		log.Fatal(err)
	}

}
