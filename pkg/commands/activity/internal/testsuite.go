package internal

import (
	"crypto/rand"
	"math/big"
	"strconv"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/tidwall/gjson"

	"github.com/bzimmer/gravl/pkg/commands/internal"
)

const N = 102

// ActivityTestSuite for testing activity services
type ActivityTestSuite struct {
	suite.Suite
	Name       string
	Encodings  []string
	SkipRoutes bool
}

func random(n int) int {
	b, err := rand.Int(rand.Reader, big.NewInt(int64(n)))
	if err != nil {
		panic(err)
	}
	return int(b.Int64())
}

func (s *ActivityTestSuite) SetupSuite() {
	if testing.Short() {
		s.T().Skip("skipping integration test suite")
		return
	}
}

func (s *ActivityTestSuite) BeforeTest(suiteName, testName string) {
	s.T().Parallel()
}

func (s *ActivityTestSuite) TestAthlete() {
	a := s.Assert()

	c := internal.Gravl("-c", s.Name, "athlete")
	<-c.Start()
	a.True(c.Success())
}

func (s *ActivityTestSuite) TestRoute() {
	if s.SkipRoutes {
		s.T().Skip("skipping routes test")
	}
	a := s.Assert()

	c := internal.Gravl("-c", s.Name, "routes", "-N", strconv.FormatInt(N, 10))
	<-c.Start()
	a.True(c.Success())
}

func (s *ActivityTestSuite) TestActivity() {
	a := s.Assert()

	c := internal.Gravl("--timeout", "30s", "-c", s.Name, "activities", "-N", strconv.FormatInt(N, 10))
	<-c.Start()
	a.True(c.Success())

	var i int
	var randomID = random(len(c.Status().Stdout))
	gjson.ForEachLine(c.Stdout(), func(res gjson.Result) bool {
		id := gjson.Get(res.String(), "id").Int()
		a.Greater(id, int64(0))
		if i == randomID {
			idS := strconv.FormatInt(id, 10)
			c = internal.Gravl("-c", s.Name, "activity", idS)
			<-c.Start()
			a.True(c.Success())
			for i := range s.Encodings {
				c = internal.Gravl("-e", s.Encodings[i], s.Name, "activity", idS)
				<-c.Start()
				if !c.Success() {
					a.FailNowf("failed encoding", s.Encodings[i])
				}
			}
			c = internal.Gravl("-c", s.Name, "activity", idS)
			<-c.Start()
			a.True(c.Success())
		}
		i++
		return true
	})
	a.Equal(N, i)
}
