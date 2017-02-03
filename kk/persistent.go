package kk

import (
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"runtime"
)

const magic = 0x15315320
const ver = 0x00000100

type Persistent struct {
	Magic  uint32
	Ver    uint32
	Height int
	Width  int
	Score  int
	Field  []int
}

func saveFname() string {
	if runtime.GOOS == "android" {
		// I'm quite sure this is not the right way to do
		// this, but I can't find any documentation that says
		// anything better.
		return "/storage/emulated/0/Android/data/org.blahonga.kk/files/state"
	}
	if home := os.Getenv("HOME"); home != "" {
		return fmt.Sprintf("%s/.kkstate", home)
	}
	return ""
}

// Failing to save or restore isn't a fatal error.
func (s *State) save() {
	fname := saveFname()
	if s.saved == nil || fname == "" {
		return
	}
	h, w := s.saved.H(), s.saved.W()
	p := Persistent{
		Magic:  magic,
		Ver:    ver,
		Height: h,
		Width:  w,
		Score:  s.saved.score,
		Field:  make([]int, h*w),
	}
	for i := range p.Field {
		p.Field[i] = s.saved.f[i/h][i%h]
	}
	f, err := os.Create(fname)
	if err != nil {
		log.Print("OpenFile ", fname, err)
		return
	}
	defer f.Close()
	enc := gob.NewEncoder(f)
	if err := enc.Encode(&p); err != nil {
		log.Print("gob.Encode ", err)
		return
	}
}

func (s *State) restore() {
	fname := saveFname()
	if fname == "" {
		return
	}
	f, err := os.Open(fname)
	if err != nil {
		return
	}
	defer f.Close()
	p := Persistent{}
	dec := gob.NewDecoder(f)
	if err := dec.Decode(&p); err != nil {
		log.Print("gob.Decode ", err)
	}
	if p.Magic != magic {
		log.Printf("Magic mismatch: %x != %x", p.Magic, magic)
		return
	}
	if p.Ver != ver {
		log.Printf("Version mismatch: %x != %x. Don't be from the future.", p.Ver, ver)
		return
	}
	h, w := s.f.H(), s.f.W()
	if p.Height != h || p.Width != w {
		log.Printf("Field size mismatch: %d != %d || %d != %d. Don't be from the future.", p.Height, h, p.Width, w)
		return
	}
	s.saved = &Field{
		score: p.Score,
	}
	for i := range p.Field {
		s.saved.f[i/h][i%h] = p.Field[i]
	}
}
