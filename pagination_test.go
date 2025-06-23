package coze

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestData 测试数据结构
type TestData struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// mockDataSource 模拟数据源
type mockDataSource struct {
	data  []*TestData
	total int
}

// newMockDataSource 创建模拟数据源
func newMockDataSource(total int) *mockDataSource {
	data := make([]*TestData, total)
	for i := 0; i < total; i++ {
		data[i] = &TestData{
			ID:   i + 1,
			Name: fmt.Sprintf("test-%d", i+1),
		}
	}
	return &mockDataSource{
		data:  data,
		total: total,
	}
}

// getNumberPageData 获取基于页码的分页数据
func (m *mockDataSource) getNumberPageData(request *pageRequest) (*pageResponse[TestData], error) {
	pageSize := request.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	startIndex := (request.PageNum - 1) * pageSize
	if startIndex >= len(m.data) {
		return &pageResponse[TestData]{
			HasMore: false,
			Total:   len(m.data),
			Data:    []*TestData{},
		}, nil
	}

	endIndex := startIndex + pageSize
	if endIndex > len(m.data) {
		endIndex = len(m.data)
	}

	return &pageResponse[TestData]{
		HasMore: endIndex < len(m.data),
		Total:   len(m.data),
		Data:    m.data[startIndex:endIndex],
	}, nil
}

// getTokenPageData 获取基于令牌的分页数据
func (m *mockDataSource) getTokenPageData(request *pageRequest) (*pageResponse[TestData], error) {
	pageSize := request.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	startIndex := 0
	if request.PageToken != "" {
		// 简单模拟：令牌就是上一页最后一条数据的 ID
		for i, item := range m.data {
			if fmt.Sprintf("%d", item.ID) == request.PageToken {
				startIndex = i + 1
				break
			}
		}
	}

	if startIndex >= len(m.data) {
		return &pageResponse[TestData]{
			HasMore: false,
			Total:   len(m.data),
			Data:    []*TestData{},
		}, nil
	}

	endIndex := startIndex + pageSize
	if endIndex > len(m.data) {
		endIndex = len(m.data)
	}

	var nextID string
	if endIndex < len(m.data) {
		nextID = fmt.Sprintf("%d", m.data[endIndex-1].ID)
	}

	return &pageResponse[TestData]{
		HasMore: endIndex < len(m.data),
		Total:   len(m.data),
		Data:    m.data[startIndex:endIndex],
		NextID:  nextID,
	}, nil
}

func TestNumberPaged(t *testing.T) {
	as := assert.New(t)
	// 创建模拟数据源
	total := 25
	mockSource := newMockDataSource(total) // 总共25条数据
	pageSize := 10

	// 创建基于页码的分页器
	pager, err := NewNumberPaged[TestData](mockSource.getNumberPageData, pageSize, 0)
	as.Nil(err)
	as.NotNil(pager)

	// 测试迭代器
	t.Run("Iterator", func(t *testing.T) {
		count := 0
		for pager.Next() {
			count++
			item := pager.Current()
			as.Equal(count, item.ID)
			as.Equal(fmt.Sprintf("test-%d", count), item.Name)
		}
		as.Equal(total, count)
		as.False(pager.HasMore())
		as.Equal(25, pager.Total())
		as.Nil(pager.Err())
	})

	t.Run("manual fetch next page", func(t *testing.T) {
		count := 0
		hasMore := true
		currentPage := 1
		for hasMore {
			pager, err := NewNumberPaged[TestData](mockSource.getNumberPageData, pageSize, currentPage)
			as.Nil(err)
			hasMore = pager.HasMore()
			count += len(pager.Items())
			currentPage++
		}
		as.Equal(total, count)
		as.False(pager.HasMore())
		as.Equal(total, pager.Total())
		as.Nil(pager.Err())
	})
}

func TestTokenPaged(t *testing.T) {
	as := assert.New(t)
	total := 25
	mockSource := newMockDataSource(total) // 总共25条数据
	pageSize := 10

	pager, err := NewLastIDPaged[TestData](mockSource.getTokenPageData, pageSize, nil)
	as.Nil(err)
	as.NotNil(pager)

	t.Run("iterator", func(t *testing.T) {
		count := 0
		for pager.Next() {
			count++
			item := pager.Current()
			as.Equal(count, item.ID)
			as.Equal(fmt.Sprintf("test-%d", count), item.Name)
		}
		as.Equal(total, count)
		as.False(pager.HasMore())
		as.Nil(pager.Err())
	})

	t.Run("manual fetch next page", func(t *testing.T) {
		count := 0
		hasMore := true
		var nextID *string
		for hasMore {
			pager, err := NewLastIDPaged[TestData](mockSource.getTokenPageData, pageSize, nextID)
			assert.Nil(t, err)
			hasMore = pager.HasMore()
			count += len(pager.Items())
			for _, item := range pager.Items() {
				if item != nil {
					nextID = ptr(strconv.Itoa(item.ID))
				}
			}
		}
		as.Equal(total, count)
		as.False(pager.HasMore())
		as.Nil(pager.Err())
	})
}

func TestPagerError(t *testing.T) {
	as := assert.New(t)
	// 测试错误情况
	errorFetcher := func(request *pageRequest) (*pageResponse[TestData], error) {
		return nil, fmt.Errorf("mock error")
	}

	// 测试基于页码的分页器错误处理
	t.Run("NumberPaged Error", func(t *testing.T) {
		pager, err := NewNumberPaged[TestData](errorFetcher, 10, 1)
		as.Error(err)
		as.Nil(pager)
	})

	// 测试基于令牌的分页器错误处理
	t.Run("TokenPaged Error", func(t *testing.T) {
		pager, err := NewLastIDPaged[TestData](errorFetcher, 10, nil)
		as.Error(err)
		as.Nil(pager)
	})
}

func TestEmptyPage(t *testing.T) {
	as := assert.New(t)
	// 创建空数据源
	emptySource := newMockDataSource(0)

	// 测试基于页码的空分页
	t.Run("Empty NumberPaged", func(t *testing.T) {
		pager, err := NewNumberPaged[TestData](emptySource.getNumberPageData, 10, 1)
		as.Nil(err)
		as.NotNil(pager)
		as.False(pager.Next())
		as.Equal(0, pager.Total())
		as.False(pager.HasMore())
		as.Nil(pager.Err())
	})

	// 测试基于令牌的空分页
	t.Run("Empty TokenPaged", func(t *testing.T) {
		pager, err := NewLastIDPaged[TestData](emptySource.getTokenPageData, 10, nil)
		as.Nil(err)
		as.NotNil(pager)
		as.False(pager.Next())
		as.False(pager.HasMore())
		as.Nil(pager.Err())
	})
}

func TestInvalidPageSize(t *testing.T) {
	as := assert.New(t)
	mockSource := newMockDataSource(25)

	// 测试基于页码的无效页大小
	t.Run("Invalid PageSize NumberPaged", func(t *testing.T) {
		pager, err := NewNumberPaged[TestData](mockSource.getNumberPageData, 0, 1)
		as.Nil(err)
		as.NotNil(pager)
		as.True(pager.Next())
		as.Equal(25, pager.Total())
	})

	// 测试基于令牌的无效页大小
	t.Run("Invalid PageSize TokenPaged", func(t *testing.T) {
		pager, err := NewLastIDPaged[TestData](mockSource.getTokenPageData, 0, nil)
		as.Nil(err)
		as.NotNil(pager)
		as.True(pager.Next())
	})
}
