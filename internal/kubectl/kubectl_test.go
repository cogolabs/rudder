package kubectl

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/cogolabs/rudder/internal/config"
	"github.com/stretchr/testify/suite"
	"gopkg.in/h2non/gock.v1"
)

const (
	testVersion = "v1.13.2"
	testBinary  = "kubectl binary"
	testYAML    = "../../test/k8s/deploy.yml"
)

type KubectlTestSuite struct {
	suite.Suite
	origYAML []byte
}

func (suite *KubectlTestSuite) SetupSuite() {
	require := suite.Require()

	b, err := ioutil.ReadFile(testYAML)
	require.NoError(err)
	suite.origYAML = b
}

func (suite *KubectlTestSuite) TearDownTest() {
	require := suite.Require()

	kubectlPath = "./kubectl"
	err := Uninstall()
	require.NoError(err)

	ioutil.WriteFile(testYAML, suite.origYAML, 0644)
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
	err := ApplyDir(buf, "../../test/k8s", "mytag", "./my/kube/config")
	require.NoError(err)
	expected := `echo apply -f ../../test/k8s/deploy.yml --kubeconfig=./my/kube/config
apply -f ../../test/k8s/deploy.yml --kubeconfig=./my/kube/config
`
	assert.Equal(expected, buf.String())
}

func (suite *KubectlTestSuite) TestSubTag() {
	assert := suite.Assert()
	require := suite.Require()

	err := subTag(testYAML, "my_tag")
	require.NoError(err)
	b, err := ioutil.ReadFile(testYAML)
	require.NoError(err)
	assert.True(strings.Contains(string(b), "my_tag"))
	assert.False(strings.Contains(string(b), imageTagPlaceholder))
	err = unstashFile(testYAML)
	require.NoError(err)
}

func (suite *KubectlTestSuite) TestWaitForRollouts() {
	assert := suite.Assert()
	require := suite.Require()
	kubectlPath = "echo"

	buf := new(bytes.Buffer)
	err := WaitForRollouts(buf, config.Deployment{KubeNamespace: "mys", KubeDeployments: []string{"myproj-dply"}})
	require.NoError(err)
	expected := `Waiting for myproj-dply in namespace mys to rollout...
echo rollout status -n mys myproj-dply
rollout status -n mys myproj-dply
`
	assert.Equal(expected, buf.String())
}

func TestKubectlTestSuite(t *testing.T) {
	tests := new(KubectlTestSuite)
	suite.Run(t, tests)
}
