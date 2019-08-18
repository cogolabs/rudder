package docker

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/cogolabs/rudder/internal/config"
	"github.com/stretchr/testify/suite"
	"gopkg.in/h2non/gock.v1"
)

const (
	testRegistry = "https://registry.server.net"
	testImage    = "myproj/api"
	testTag      = "v3.1"
)

var testConfig = config.Config{
	Containers: []config.Container{{
		Registry: testRegistry,
		Image:    testImage,
		Timeout:  10 * time.Millisecond,
	}},
}

type DockerTestSuite struct {
	suite.Suite
}

func (suite *DockerTestSuite) SetupSuite() {
	tickerInterval = time.Millisecond
}

func (suite *DockerTestSuite) TearDownTest() {
	gock.Off()
}

func (suite *DockerTestSuite) TestCheckImage() {
	assert := suite.Assert()
	require := suite.Require()

	gock.New(testRegistry).
		Get(fmt.Sprintf("/v1/repositories/%s/tags/%s", testImage, testTag)).
		Reply(http.StatusOK)

	ready, err := checkImage(&testConfig.Containers[0], testTag)
	require.NoError(err)
	assert.True(ready)
}

func (suite *DockerTestSuite) TestCheckImageNotReady() {
	assert := suite.Assert()
	require := suite.Require()

	gock.New(testRegistry).
		Get(fmt.Sprintf("/v1/repositories/%s/tags/%s", testImage, testTag)).
		Reply(http.StatusNotFound)

	ready, err := checkImage(&testConfig.Containers[0], testTag)
	require.NoError(err)
	assert.False(ready)
}

func (suite *DockerTestSuite) TestCheckImageRegistryError() {
	assert := suite.Assert()
	require := suite.Require()

	gock.New(testRegistry).
		Get(fmt.Sprintf("/v1/repositories/%s/tags/%s", testImage, testTag)).
		Reply(http.StatusInternalServerError)

	ready, err := checkImage(&testConfig.Containers[0], testTag)
	url := fmt.Sprintf("%s/v1/repositories/%s/tags/%s", testRegistry, testImage, testTag)
	require.EqualError(err, fmt.Sprintf("received code 500 from '%s'", url))
	assert.False(ready)
}

func (suite *DockerTestSuite) TestWaitForImage() {
	require := suite.Require()

	gock.New(testRegistry).
		Get(fmt.Sprintf("/v1/repositories/%s/tags/%s", testImage, testTag)).
		Reply(http.StatusNotFound)
	gock.New(testRegistry).
		Get(fmt.Sprintf("/v1/repositories/%s/tags/%s", testImage, testTag)).
		Reply(http.StatusOK)

	err := WaitForImages(&testConfig, testTag)
	require.NoError(err)
}

func (suite *DockerTestSuite) TestWaitForImageTimeout() {
	require := suite.Require()

	gock.New(testRegistry).
		Get(fmt.Sprintf("/v1/repositories/%s/tags/%s", testImage, testTag)).
		Persist().
		Reply(http.StatusNotFound)

	err := WaitForImages(&testConfig, testTag)
	require.Equal(ErrTimeout, err)
}

func (suite *DockerTestSuite) TestWaitForImageError() {
	require := suite.Require()

	gock.New(testRegistry).
		Get(fmt.Sprintf("/v1/repositories/%s/tags/%s", testImage, testTag)).
		Reply(http.StatusInternalServerError)

	err := WaitForImages(&testConfig, testTag)
	url := fmt.Sprintf("%s/v1/repositories/%s/tags/%s", testRegistry, testImage, testTag)
	require.EqualError(err, fmt.Sprintf("received code 500 from '%s'", url))
}

func TestDockerTestSuite(t *testing.T) {
	tests := new(DockerTestSuite)
	suite.Run(t, tests)
}
