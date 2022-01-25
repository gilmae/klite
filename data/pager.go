package data

const PageSize = 4096

type Page [PageSize]byte

type Pager interface {
	NextPage() Page
}
