package docker

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/ryantking/rudder/internal/config"
	"github.com/stretchr/testify/suite"
	"gopkg.in/h2non/gock.v1"
)

const testTag = "v3.1"

var testConfig = config.Config{
	DockerRegistry: "https://registry.server.net",
	DockerImage:    "myproj/api",
	DockerTimeout:  10 * time.Millisecond,
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

	gock.New(testConfig.DockerRegistry).
		Get(fmt.Sprintf("/v1/repositories/%s/tags/%s", testConfig.DockerImage, testTag)).
		Reply(http.StatusOK)

	ready, err := checkImage(&testConfig, testTag)
	require.NoError(err)
	assert.True(ready)
}

func (suite *DockerTestSuite) TestCheckImageNotReady() {
	assert := suite.Assert()
	require := suite.Require()

	gock.New(testConfig.DockerRegistry).
		Get(fmt.Sprintf("/v1/repositories/%s/tags/%s", testConfig.DockerImage, testTag)).
		Reply(http.StatusNotFound)

	ready, err := checkImage(&testConfig, testTag)
	require.NoError(err)
	assert.False(ready)
}

func (suite *DockerTestSuite) TestCheckImageRegistryError() {
	assert := suite.Assert()
	require := suite.Require()

	gock.New(testConfig.DockerRegistry).
		Get(fmt.Sprintf("/v1/repositories/%s/tags/%s", testConfig.DockerImage, testTag)).
		Reply(http.StatusInternalServerError)

	ready, err := checkImage(&testConfig, testTag)
	url := fmt.Sprintf("%s/v1/repositories/%s/tags/%s", testConfig.DockerRegistry, testConfig.DockerImage, testTag)
	require.EqualError(err, fmt.Sprintf("recieved code 500 from '%s'", url))
	assert.False(ready)
}

func (suite *DockerTestSuite) TestWaitForImage() {
	require := suite.Require()

	gock.New(testConfig.DockerRegistry).
		Get(fmt.Sprintf("/v1/repositories/%s/tags/%s", testConfig.DockerImage, testTag)).
		Reply(http.StatusNotFound)
	gock.New(testConfig.DockerRegistry).
		Get(fmt.Sprintf("/v1/repositories/%s/tags/%s", testConfig.DockerImage, testTag)).
		Reply(http.StatusOK)

	err := WaitForImage(&testConfig, testTag)
	require.NoError(err)
}

func (suite *DockerTestSuite) TestWaitForImageTimeout() {
	require := suite.Require()

	gock.New(testConfig.DockerRegistry).
		Get(fmt.Sprintf("/v1/repositories/%s/tags/%s", testConfig.DockerImage, testTag)).
		Persist().
		Reply(http.StatusNotFound)

	err := WaitForImage(&testConfig, testTag)
	require.Equal(ErrTimeout, err)
}

func (suite *DockerTestSuite) TestWaitForImageError() {
	require := suite.Require()

	gock.New(testConfig.DockerRegistry).
		Get(fmt.Sprintf("/v1/repositories/%s/tags/%s", testConfig.DockerImage, testTag)).
		Reply(http.StatusInternalServerError)

	err := WaitForImage(&testConfig, testTag)
	url := fmt.Sprintf("%s/v1/repositories/%s/tags/%s", testConfig.DockerRegistry, testConfig.DockerImage, testTag)
	require.EqualError(err, fmt.Sprintf("recieved code 500 from '%s'", url))
}

func TestDockerTestSuite(t *testing.T) {
	tests := new(DockerTestSuite)
	suite.Run(t, tests)
}
