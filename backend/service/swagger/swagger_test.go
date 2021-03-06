package swagger_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ekspand/trusty/backend/service/swagger"
	"github.com/ekspand/trusty/internal/appcontainer"
	"github.com/ekspand/trusty/internal/config"
	"github.com/ekspand/trusty/pkg/gserver"
	"github.com/ekspand/trusty/tests/testutils"
	"github.com/go-phorce/dolly/xhttp/header"
	"github.com/go-phorce/dolly/xhttp/retriable"
	"github.com/go-phorce/dolly/xlog"
	"github.com/juju/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	trustyServer *gserver.Server
	httpAddr     string
	httpsAddr    string
)

const (
	projFolder = "../../../"
)

// serviceFactories provides map of trustyserver.ServiceFactory
var serviceFactories = map[string]gserver.ServiceFactory{
	swagger.ServiceName: swagger.Factory,
}

func TestMain(m *testing.M) {
	var err error
	xlog.SetPackageLogLevel("github.com/go-phorce/dolly/xhttp", "retriable", xlog.DEBUG)

	httpsAddr = testutils.CreateURLs("https", "")
	httpAddr = testutils.CreateURLs("http", "")

	devcfg, err := testutils.LoadConfig(projFolder, "UNIT_TEST")
	if err != nil {
		panic(err.Error())
	}
	cfg := &config.HTTPServer{
		ListenURLs: []string{httpsAddr, httpAddr},
		ServerTLS: &config.TLSInfo{
			CertFile:      "/tmp/trusty/certs/trusty_dev_peer_wfe.pem",
			KeyFile:       "/tmp/trusty/certs/trusty_dev_peer_wfe-key.pem",
			TrustedCAFile: "/tmp/trusty/certs/trusty_dev_root_ca.pem",
		},
		Services: []string{swagger.ServiceName},
		Swagger:  devcfg.HTTPServers["cis"].Swagger,
	}

	container := appcontainer.NewBuilder().
		WithAuditor(nil).
		WithCrypto(nil).
		WithJwtParser(nil).
		WithDiscovery(appcontainer.NewDiscovery()).
		Container()

	trustyServer, err = gserver.Start("SwaggerTest", cfg, container, serviceFactories)
	if err != nil || trustyServer == nil {
		panic(errors.Trace(err))
	}

	// TODO: channel for <-trustyServer.ServerReady()

	// Run the tests
	rc := m.Run()

	// cleanup
	trustyServer.Close()

	os.Exit(rc)
}

func TestSwagger(t *testing.T) {
	client := retriable.New()

	w := httptest.NewRecorder()
	hdr, status, err := client.Request(context.Background(),
		http.MethodGet,
		[]string{httpAddr},
		"/v1/swagger/status",
		nil,
		w)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, status)
	assert.Contains(t, hdr.Get(header.ContentType), header.ApplicationJSON)

	w = httptest.NewRecorder()
	_, status, err = client.Request(context.Background(),
		http.MethodGet,
		[]string{httpAddr},
		"/v1/swagger/notfound",
		nil,
		w)
	require.Error(t, err)
	assert.Equal(t, http.StatusNotFound, status)
}
