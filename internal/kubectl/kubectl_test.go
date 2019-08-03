package kubectl

import (
	"bytes"
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

	kubectlPath = "./kubectl"
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

func (suite *KubectlTestSuite) TestApplyDir() {
	assert := suite.Assert()
	require := suite.Require()
	kubectlPath = "echo"

	buf := new(bytes.Buffer)
	err := ApplyDir(buf, "../../test/k8s", "./my/kube/config")
	require.NoError(err)
	expected := `echo apply -f ../../test/k8s/deploy.yml --kubeconfig=./my/kube/config
apply -f ../../test/k8s/deploy.yml --kubeconfig=./my/kube/config
`
	assert.Equal(expected, buf.String())
}

func TestKubectlTestSuite(t *testing.T) {
	tests := new(KubectlTestSuite)
	suite.Run(t, tests)
}
