package data

import (
	"fmt"
	"os"
)

const PageSize = 4096
const MAXPAGES = uint32(1000)

type Page []byte

type Pager interface {
	Page(page uint32) (*Page, error)
	GetNextUnusedPageNum() uint32
	Close()
	Flush() error
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

func (mp *MemoryPager) Close() {}

func (mp *MemoryPager) Flush() error { return nil }

type FilePager struct {
	fileDescriptor *os.File
	fileLength     int64
	pages          [MAXPAGES]Page
	NumPages       uint32
}

func NewFilePager(filename string) (*FilePager, error) {
	p := FilePager{}
	var file *os.File
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		file, err = os.Create(filename)
		if err != nil {
			return nil, err
		}
	} else {
		file, err = os.OpenFile(filename, os.O_RDWR, os.ModeExclusive)

		if err != nil {
			return nil, err
		}
	}
	p.fileDescriptor = file

	p.fileLength, _ = p.fileDescriptor.Seek(0, 2)
	for i, _ := range p.pages {
		p.pages[i] = nil
	}
	p.NumPages = uint32(p.fileLength) / uint32(PageSize)

	if p.fileLength%int64(PageSize) != 0 {
		return nil, fmt.Errorf("db file is not a whole number of pages, corrupt file")
	}
	return &p, nil
}

func (p *FilePager) Close() {
	p.fileDescriptor.Close()
}

func (p *FilePager) Flush() error {
	for i, page := range p.pages {

		if page == nil {
			continue
		}

		position := int64(i * int(PageSize))
		_, err := p.fileDescriptor.Seek(position, 0)
		if err != nil {
			return err
		}

		bytesWritten, err := p.fileDescriptor.Write(page)

		if err != nil {
			return err
		}

		if bytesWritten != int(PageSize) {
			return fmt.Errorf("incorrect number of bytes written: %d", bytesWritten)
		}
	}
	return nil
}

func (p *FilePager) GetNextUnusedPageNum() uint32 {
	return p.NumPages
}

func (p *FilePager) Page(pageNum uint32) (*Page, error) {
	if pageNum > MAXPAGES {
		return nil, fmt.Errorf("pageNum out of bounds, max pages: %d", MAXPAGES)
	}

	if p.pages[pageNum] == nil {
		p.pages[pageNum] = make([]byte, PageSize)
		num_pages := p.fileLength / int64(PageSize)

		if p.fileLength%int64(PageSize) != 0 {
			num_pages += 1
		}

		if int64(pageNum) < num_pages {
			p.fileDescriptor.Seek(int64(pageNum)*int64(PageSize), 0)

			bytesRead, err := p.fileDescriptor.Read(p.pages[pageNum])
			if err != nil {
				return nil, err
			}
			if bytesRead == -1 {
				return nil, fmt.Errorf("error reading file")
			}

		}

	}

	if pageNum >= p.NumPages {
		p.NumPages = pageNum + 1
	}
	page := p.pages[pageNum]
	return &page, nil
}
