package coze

type BasePaged[T any] interface {
	Err() error
	Items() []*T
	Current() *T
	Next() bool
	HasMore() bool
}

type NumberPaged[T any] interface {
	BasePaged[T]
	Total() int
}

type LastIDPaged[T any] interface {
	BasePaged[T]
	GetLastID() string
}

type pageRequest struct {
	PageToken string `json:"page_token,omitempty"`
	PageNum   int    `json:"page_num,omitempty"`
	PageSize  int    `json:"page_size,omitempty"`
}

type pageResponse[T any] struct {
	HasMore bool   `json:"has_more"`
	Total   int    `json:"total"`
	Data    []*T   `json:"data"`
	LastID  string `json:"last_id,omitempty"`
	NextID  string `json:"next_id,omitempty"`
	LogID   string `json:"log_id,omitempty"`
}

type basePager[T any] struct {
	pageFetcher    PageFetcher[T]
	pageSize       int
	currentPage    *pageResponse[T]
	currentIndex   int
	currentPageNum int
	cur            *T
	err            error
}

func (p *basePager[T]) Err() error {
	return p.err
}

func (p *basePager[T]) Items() []*T {
	return ptrValue(p.currentPage).Data
}

func (p *basePager[T]) Current() *T {
	return p.cur
}

func (p *basePager[T]) Total() int {
	return ptrValue(p.currentPage).Total
}

func (p *basePager[T]) HasMore() bool {
	return ptrValue(p.currentPage).HasMore
}

// PageFetcher interface
type PageFetcher[T any] func(request *pageRequest) (*pageResponse[T], error)

// NumberPaged implementation
type implNumberPaged[T any] struct {
	basePager[T]
}

func NewNumberPaged[T any](fetcher PageFetcher[T], pageSize, pageNum int) (NumberPaged[T], error) {
	if pageNum <= 0 {
		pageNum = 1
	}
	paginator := &implNumberPaged[T]{basePager: basePager[T]{pageFetcher: fetcher, pageSize: pageSize, currentPageNum: pageNum}}
	err := paginator.fetchNextPage()
	if err != nil {
		return nil, err
	}
	return paginator, nil
}

func (p *implNumberPaged[T]) fetchNextPage() error {
	request := &pageRequest{PageNum: p.currentPageNum, PageSize: p.pageSize}
	var err error
	p.currentPage, err = p.pageFetcher(request)
	if err != nil {
		return err
	}
	p.currentIndex = 0
	p.currentPageNum++
	return nil
}

func (p *implNumberPaged[T]) Next() bool {
	if p.currentIndex < len(ptrValue(p.currentPage).Data) {
		p.cur = p.currentPage.Data[p.currentIndex]
		p.currentIndex++
		return true
	}
	if p.currentPage.HasMore {
		err := p.fetchNextPage()
		if err != nil {
			p.err = err
			return false
		}
		if len(p.currentPage.Data) == 0 {
			return false
		}
		p.cur = p.currentPage.Data[p.currentIndex]
		p.currentIndex++
		return true
	}
	return false
}

// TokenPaged implementation
type implLastIDPaged[T any] struct {
	basePager[T]
	pageToken *string
}

func NewLastIDPaged[T any](fetcher PageFetcher[T], pageSize int, nextID *string) (LastIDPaged[T], error) {
	paginator := &implLastIDPaged[T]{basePager: basePager[T]{pageFetcher: fetcher, pageSize: pageSize}, pageToken: nextID}
	err := paginator.fetchNextPage()
	if err != nil {
		return nil, err
	}
	return paginator, nil
}

func (p *implLastIDPaged[T]) fetchNextPage() error {
	request := &pageRequest{PageToken: ptrValue(p.pageToken), PageSize: p.pageSize}
	var err error
	p.currentPage, err = p.pageFetcher(request)
	if err != nil {
		return err
	}
	p.currentIndex = 0
	p.pageToken = &p.currentPage.NextID
	return nil
}

func (p *implLastIDPaged[T]) Next() bool {
	if p.currentIndex < len(ptrValue(p.currentPage).Data) {
		p.cur = p.currentPage.Data[p.currentIndex]
		p.currentIndex++
		return true
	}
	if p.currentPage.HasMore {
		err := p.fetchNextPage()
		if err != nil {
			p.err = err
			return false
		}
		if len(p.currentPage.Data) == 0 {
			return false
		}
		p.cur = p.currentPage.Data[p.currentIndex]
		p.currentIndex++
		return true
	}
	return false
}

func (p *implLastIDPaged[T]) GetLastID() string {
	return p.currentPage.LastID
}
