package data

import "fmt"

const PageSize = 4096
const MAXPAGES = uint32(100)

type Page []byte

type Pager interface {
	Page(page uint32) (*Page, error)
	GetNextUnusedPageNum() uint32
}

type MemoryPager struct {
	pages    [MAXPAGES]Page
	nextPage uint32
}

func (mp *MemoryPager) GetNextUnusedPageNum() uint32 {
	return mp.nextPage
}

func (mp *MemoryPager) Page(page uint32) (*Page, error) {
	if page > MAXPAGES {
		return nil, fmt.Errorf("page out of bounds, max pages: %d", MAXPAGES)
	}

	if mp.pages[page] == nil {
		mp.pages[page] = make([]byte, PageSize)
	}

	if page >= mp.nextPage {
		mp.nextPage = page + 1
	}

	return &mp.pages[page], nil
}
