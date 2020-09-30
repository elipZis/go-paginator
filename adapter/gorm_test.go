package adapter_test

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"github.com/vcraescu/go-paginator"
	"github.com/vcraescu/go-paginator/adapter"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

type (
	Post struct {
		ID     uint `gorm:"primary_key"`
		Number int
	}

	GORMAdapterTestSuite struct {
		suite.Suite
		db *gorm.DB
	}
)

func (suite *GORMAdapterTestSuite) SetupTest() {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("setup test: %s", err))
	}

	suite.db = db
	suite.db.AutoMigrate(&Post{})

	for i := 1; i <= 100; i++ {
		p := Post{
			Number: i,
		}

		suite.db.Save(&p)
	}
}

func (suite *GORMAdapterTestSuite) TearDownTest() {
	rawDB, _ := suite.db.DB()
	if err := rawDB.Close(); err != nil {
		panic(fmt.Errorf("tear down test: %s", err))
	}
}

func (suite *GORMAdapterTestSuite) TestFirstPage() {
	q := suite.db.Model(Post{})
	p := paginator.New(adapter.NewGORMAdapter(q), 10)

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

func (suite *GORMAdapterTestSuite) TestLastPage() {
	q := suite.db.Model(Post{})
	p := paginator.New(adapter.NewGORMAdapter(q), 10)

	p.SetPage(10)

	hn, _ := p.HasNext()
	suite.False(hn)

	hp, _ := p.HasPrev()
	suite.True(hp)
}

func (suite *GORMAdapterTestSuite) TestOutOfRangeCurrentPage() {
	q := suite.db.Model(Post{})
	p := paginator.New(adapter.NewGORMAdapter(q), 10)

	var posts []Post
	p.SetPage(11)
	err := p.Results(&posts)
	suite.NoError(err)

	page, _ := p.Page()
	suite.Equal(10, page)

	posts = make([]Post, 0)
	p.SetPage(-4)

	hn, _ := p.HasNext()
	suite.True(hn)

	hp, _ := p.HasPrev()
	suite.False(hp)

	hpages, _ := p.HasPages()
	suite.True(hpages)

	err = p.Results(&posts)
	suite.NoError(err)
	suite.Len(posts, 10)
}

func (suite *GORMAdapterTestSuite) TestCurrentPageResults() {
	q := suite.db.Model(Post{})
	p := paginator.New(adapter.NewGORMAdapter(q), 10)

	var posts []Post
	p.SetPage(6)
	err := p.Results(&posts)
	suite.NoError(err)

	suite.Len(posts, 10)
	for i, post := range posts {
		page, err := p.Page()
		if !suite.NoError(err) {
			return
		}
		fmt.Println(p.Nums())
		suite.Equal((page-1)*10+i+1, post.Number)
	}
}

func TestGORMAdapterTestSuite(t *testing.T) {
	suite.Run(t, new(GORMAdapterTestSuite))
}
