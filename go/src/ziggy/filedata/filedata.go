package filedata

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
)

const (
	structTagLabel = "fd"
	uniqueKey      = "key" // This is what will be used as the basis to replace/ or append
)

// Mgr manages array's of json encoded metadata in a file.
// It assumes that all data put into that file is serializable to the same type
// using the builtin json pkg
type Mgr interface {
	// Put replaces the existing metadata record (with a specified primary key)
	Put(interface{}) error
	// Append just adds a new record to the bottom of the list
	Append(interface{}) error
}

func NewMgr(loc string) Mgr {
	return &mgr{loc: loc}
}

type mgr struct {
	loc string
}

func (m *mgr) replaceOrAppend(v interface{}, appendOnly bool) error {
	vType := reflect.TypeOf(v)
	vVal := reflect.ValueOf(v)
	if vType.Kind() != reflect.Ptr && vType.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("interface is not a ptr to a struct")
	}

	keyIndex := getKeyIndex(vType)

	key := vVal.Elem().Field(keyIndex).String()

	// Create a pointer to a slice of the type
	mdVal := reflect.New(reflect.SliceOf(vType))

	if err := readData(m.loc, mdVal); err != nil {
		return err
	}

	if !appendOnly {
		var replaced bool
		for i := 0; i < mdVal.Elem().Len(); i++ {
			keyAtI := mdVal.Elem().Index(i).Elem().Field(keyIndex).String()
			if keyAtI == key {
				mdVal.Elem().Index(i).Set(vVal)
				replaced = true
				break
			}
		}

		if !replaced {
			mdVal.Elem().Set(reflect.Append(mdVal.Elem(), vVal))
		}

	} else {
		mdVal.Elem().Set(reflect.Append(mdVal.Elem(), vVal))
	}

	bkPath := fmt.Sprintf("%s.bk", m.loc)
	if err := writeFile(bkPath, mdVal.Interface()); err != nil {
		return err
	}

	if err := writeFile(m.loc, mdVal.Interface()); err != nil {
		return err
	}

	os.Remove(bkPath)

	return nil
}

func (m *mgr) Put(v interface{}) error {
	return m.replaceOrAppend(v, false)
}

func (m *mgr) Append(v interface{}) error {
	return m.replaceOrAppend(v, true)
}

func getKeyIndex(vType reflect.Type) int {
	keyIndex := 0
	for i := 0; i < vType.Elem().NumField(); i++ {
		f := vType.Elem().Field(i)
		tagVal := f.Tag.Get(structTagLabel)
		if tagVal == "" {
			continue
		}

		for _, t := range strings.Split(tagVal, ",") {
			if t == uniqueKey {
				keyIndex = i
			}
		}
	}
	return keyIndex
}

func readData(loc string, mdVal reflect.Value) error {
	f, err := os.Open(loc)
	if err != nil {
		if os.IsNotExist(err) {
			f, err = os.Create(loc)
			f.Write([]byte("[]"))
			if err != nil {
				return err
			}
		}
		return err
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(mdVal.Interface()); err != nil {
		return err
	}

	return nil
}

func writeFile(loc string, v interface{}) error {
	wF, err := os.Create(loc)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(wF)
	encoder.SetIndent("", " ")

	if err := encoder.Encode(v); err != nil {
		return nil
	}

	return nil
}
