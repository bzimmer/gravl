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

const numberOfActivities = 102

// ActivityTestSuite for testing activity services
type ActivityTestSuite struct {
	suite.Suite
	Name          string
	Encodings     []string
	MaxActivities int64
	Routes        bool
	Upload        bool
	StreamSets    bool
}

func random(n int) int {
	b, err := rand.Int(rand.Reader, big.NewInt(int64(n)))
	if err != nil {
		panic(err)
	}
	return int(b.Int64())
}

func (s *ActivityTestSuite) N() int64 {
	if s.MaxActivities > 0 && s.MaxActivities < numberOfActivities {
		return s.MaxActivities
	}
	return numberOfActivities
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

func (s *ActivityTestSuite) TestRoutes() {
	// @todo(bzimmer)
	// This test needs to be improved, it currently only confirms the subcommand exists
	if !s.Routes {
		s.T().Logf("skipping routes for %s", s.Name)
		return
	}
	a := s.Assert()
	c := internal.Gravl("-c", s.Name, "routes", "-N", strconv.FormatInt(s.N(), 10))
	<-c.Start()
	a.True(c.Success())
}

func (s *ActivityTestSuite) TestStreamSets() {
	if !s.StreamSets {
		s.T().Logf("skipping streamsets for %s", s.Name)
		return
	}
	a := s.Assert()
	c := internal.Gravl("-c", s.Name, "streamsets")
	<-c.Start()
	a.True(c.Success())
	for _, x := range []string{"distance", "heartrate"} {
		y := gjson.Get(c.Stdout(), x).String()
		a.NotEmpty(y)
	}
}

func (s *ActivityTestSuite) TestUpload() {
	if !s.Upload {
		s.T().Logf("skipping upload for %s", s.Name)
		return
	}
	a := s.Assert()
	c := internal.Gravl("--timeout", "30s", "-c", s.Name, "upload", "-n")
	<-c.Start()
	a.True(c.Success())
}

func (s *ActivityTestSuite) TestActivity() {
	a := s.Assert()

	n := s.N()
	c := internal.Gravl("--timeout", "30s", "-c", s.Name, "activities", "-N", strconv.FormatInt(n, 10))
	<-c.Start()
	a.True(c.Success())

	lines := c.Status().Stdout
	a.Equal(n, int64(len(lines)))

	randomID := random(len(c.Status().Stdout))
	for i := 0; i < len(lines); i++ {
		res := gjson.Parse(lines[i])
		id := gjson.Get(res.String(), "id").Int()
		a.Greater(id, int64(0))
		if i == randomID {
			idS := strconv.FormatInt(id, 10)
			c = internal.Gravl("-c", s.Name, "activity", idS)
			<-c.Start()
			a.True(c.Success())
			for j := range s.Encodings {
				c = internal.Gravl("-e", s.Encodings[j], s.Name, "activity", idS)
				<-c.Start()
				if !c.Success() {
					a.FailNowf("failed encoding", s.Encodings[j])
				}
			}
			c = internal.Gravl("-c", s.Name, "activity", idS)
			<-c.Start()
			a.True(c.Success())
		}
	}
}
