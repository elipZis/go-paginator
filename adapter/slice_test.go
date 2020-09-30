package adapter_test

import (
	"github.com/stretchr/testify/suite"
	"github.com/vcraescu/go-paginator"
	"github.com/vcraescu/go-paginator/adapter"
	"testing"
)

type ArrayAdapterTestSuite struct {
	suite.Suite
	data []int
}

func (suite *ArrayAdapterTestSuite) SetupTest() {
	suite.data = make([]int, 100)
	for i := 1; i <= 100; i++ {
		suite.data[i-1] = i
	}
}

func (suite *ArrayAdapterTestSuite) TestFirstPage() {
	p := paginator.New(adapter.NewSliceAdapter(suite.data), 10)

	pn, _ := p.PageNums()
	suite.Equal(10, pn)

	page, _ := p.Page()
	suite.Equal(1, page)

	hn, _ := p.HasNext()
	suite.True(hn)

	hp, _ := p.HasPrev()
	suite.False(hp)

	hpages, _ := p.HasPages()
	suite.True(hpages)
}

func (suite *ArrayAdapterTestSuite) TestLastPage() {
	p := paginator.New(adapter.NewSliceAdapter(suite.data), 10)

	p.SetPage(10)

	hn, _ := p.HasNext()
	suite.False(hn)

	hp, _ := p.HasPrev()
	suite.True(hp)
}

func (suite *ArrayAdapterTestSuite) TestOutOfRangeCurrentPage() {
	p := paginator.New(adapter.NewSliceAdapter(suite.data), 10)

	var pages []int
	p.SetPage(11)
	err := p.Results(&pages)
	suite.NoError(err)

	page, _ := p.Page()
	suite.Equal(10, page)

	pages = make([]int, 0)
	p.SetPage(-4)

	hn, _ := p.HasNext()
	suite.True(hn)

	hp, _ := p.HasPrev()
	suite.False(hp)

	hpages, _ := p.HasPages()
	suite.True(hpages)

	err = p.Results(&pages)
	suite.NoError(err)
	suite.Len(pages, 10)
}

func (suite *ArrayAdapterTestSuite) TestCurrentPageResults() {
	p := paginator.New(adapter.NewSliceAdapter(suite.data), 10)

	var pages []int
	p.SetPage(6)
	err := p.Results(&pages)
	suite.NoError(err)

	suite.Len(pages, 10)
	for i, page := range pages {
		expectedPage, _ := p.Page()
		suite.Equal((expectedPage-1)*10+i+1, page)
	}
}

func TestArrayAdapterTestSuite(t *testing.T) {
	suite.Run(t, new(ArrayAdapterTestSuite))
}
