package environment

import (
	"encoding/binary"
	"unsafe"

	"github.com/gilmae/klite/data"
	"github.com/gilmae/klite/store"
)

var VERSION = []uint8{0, 9, 1}

const (
	RootPage              = 0
	IdentifierOffset      = 0
	IdentifierSize        = uint16(unsafe.Sizeof([]byte("klite")))
	VersionMajorOffset    = IdentifierSize
	VersionMajorSize      = uint16(unsafe.Sizeof(uint8(0)))
	VersionMinorOffset    = VersionMajorSize + VersionMajorOffset
	VersionMinorSize      = uint16(unsafe.Sizeof(uint8(0)))
	VersionRevisionOffset = VersionMinorOffset + VersionMinorSize
	VersionRevisionSize   = uint16(unsafe.Sizeof(uint8(0)))
	StreamPageOffset      = VersionRevisionOffset + VersionRevisionSize + 1
	StreamPageSize        = uint16(unsafe.Sizeof(uint32(0)))

	StoreHeaderSize = StreamPageOffset + StreamPageSize
)

type Environment struct {
	pager data.Pager
	page  *data.Page
}

func NewEnvironment(pager data.Pager) (*Environment, error) {
	rootPage, err := pager.Page(RootPage)
	if err != nil {
		return nil, err
	}

	return &Environment{pager: pager, page: rootPage}, nil
}

func (e *Environment) Initialise() error {
	rootPage, err := e.pager.Page(RootPage)
	e.page = rootPage
	if err != nil {
		return err
	}

	copy((*rootPage)[IdentifierOffset:IdentifierOffset+IdentifierSize], []byte("klite"))

	e.SetVersion(VERSION)
	_, streamPageNum := store.InitialiseStream(e.pager)
	e.SetStreamPage(streamPageNum)

	return nil
}

func (e *Environment) IsInitialised() bool {
	return string((*e.page)[IdentifierOffset:IdentifierOffset+IdentifierSize]) == "klite"
}

func (e *Environment) Version() []uint8 {
	version := make([]uint8, 3)
	version[0] = (*e.page)[VersionMajorOffset]
	version[1] = (*e.page)[VersionMinorOffset]
	version[2] = (*e.page)[VersionRevisionOffset]

	return version
}

func (e *Environment) SetVersion(version []uint8) {
	(*e.page)[VersionMajorOffset] = version[0]
	(*e.page)[VersionMinorOffset] = version[1]
	(*e.page)[VersionRevisionOffset] = version[2]
}

func (e *Environment) StreamPage() uint32 {
	return binary.LittleEndian.Uint32((*e.page)[StreamPageOffset : StreamPageOffset+StreamPageSize])
}

func (e *Environment) SetStreamPage(streamPage uint32) {
	binary.LittleEndian.PutUint32((*e.page)[StreamPageOffset:StreamPageOffset+StreamPageSize], streamPage)
}

func (e *Environment) GetStream() *store.Stream {
	return store.NewStream(e.pager, e.StreamPage())
}
