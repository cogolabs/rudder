package kubectl

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"testing"

	"github.com/stretchr/testify/suite"
	"gopkg.in/h2non/gock.v1"
)

const (
	testVersion = "v1.13.2"
	testBinary  = "kubectl binary"
)

type KubectlTestSuite struct {
	suite.Suite
}

func (suite *KubectlTestSuite) TearDownTest() {
	require := suite.Require()

	err := Uninstall()
	require.NoError(err)
}

func (suite *KubectlTestSuite) TestInstall() {
	assert := suite.Assert()
	require := suite.Require()

	path := fmt.Sprintf(pathBase, testVersion, runtime.GOOS, runtime.GOARCH)
	gock.New(kubectlBase).
		Get(path).
		Reply(http.StatusOK).
		BodyString(testBinary)

	err := Install(testVersion)
	require.NoError(err)
	f, err := os.Open(kubectlPath)
	require.NoError(err)
	b, err := ioutil.ReadAll(f)
	require.NoError(err)
	assert.EqualValues(testBinary, b)
}

func (suite *KubectlTestSuite) TestInstallBadResponse() {
	require := suite.Require()

	path := fmt.Sprintf(pathBase, testVersion, runtime.GOOS, runtime.GOARCH)
	gock.New(kubectlBase).
		Get(path).
		Reply(http.StatusInternalServerError).
		BodyString(testBinary)

	err := Install(testVersion)
	require.EqualError(err, "could not install kubectl, received code 500")
}

func TestKubectlTestSuite(t *testing.T) {
	tests := new(KubectlTestSuite)
	suite.Run(t, tests)
}
