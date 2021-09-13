package v1

/*
Copyright © 2021 Remi Ferrand

Contributor(s): Remi Ferrand <riton.github_at_gmail(dot)com>, 2021

This software is a computer program whose purpose is to interact with
the ENSAP API (ensap.gouv.fr).

This software is governed by the CeCILL-B license under French law and
abiding by the rules of distribution of free software.  You can  use,
modify and/ or redistribute the software under the terms of the CeCILL-B
license as circulated by CEA, CNRS and INRIA at the following URL
"http://www.cecill.info".

As a counterpart to the access to the source code and  rights to copy,
modify and redistribute granted by the license, users are provided only
with a limited warranty  and the software's author,  the holder of the
economic rights,  and the successive licensors  have only  limited
liability.

In this respect, the user's attention is drawn to the risks associated
with loading,  using,  modifying and/or developing or reproducing the
software by the user in light of its specific status of free software,
that may mean  that it is complicated to manipulate,  and  that  also
therefore means  that it is reserved for developers  and  experienced
professionals having in-depth computer knowledge. Users are therefore
encouraged to load and test the software's suitability as regards their
requirements in conditions enabling the security of their systems and/or
data to be ensured and,  more generally, to use and operate it in the
same conditions as regards security.

The fact that you are presently reading this means that you have had
knowledge of the CeCILL-B license and that you accept its terms.

*/

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/pkg/errors"
	"github.com/riton/ensap-cli/httputils"
	log "github.com/sirupsen/logrus"
)

type ensapV1APIClient struct {
	Endpoint   string
	username   string
	password   string
	cfg        EnsapV1APIClientConfig
	ctx        context.Context
	httpClient httputils.HTTPDoer
	userAgent  string
}

type EnsapV1APIClientConfig struct {
	HTTPTimeout time.Duration
	Debug       bool
}

var DefaultEnsapV1APIClientConfig = EnsapV1APIClientConfig{
	HTTPTimeout: 10 * time.Second,
	Debug:       false,
}

const (
	DefaultEnsapV1APIClientUserAgent = "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:90.0) Gecko/20100101 Firefox/90.0"
)

func NewEnsapV1APIClientWithConfig(endpoint, username, password string, cfg *EnsapV1APIClientConfig) *ensapV1APIClient {
	if cfg == nil {
		cfg = &DefaultEnsapV1APIClientConfig
	}
	return &ensapV1APIClient{
		Endpoint:  endpoint,
		username:  username,
		password:  password,
		cfg:       *cfg,
		ctx:       context.Background(),
		userAgent: DefaultEnsapV1APIClientUserAgent,
	}
}

func NewEnsapV1APIClientWithContext(ctx context.Context, endpoint, username, password string, cfg *EnsapV1APIClientConfig) *ensapV1APIClient {
	c := NewEnsapV1APIClientWithConfig(endpoint, username, password, cfg)
	c.ctx = ctx
	return c
}

func (c *ensapV1APIClient) SetUserAgent(s string) {
	c.userAgent = s
}

func (c *ensapV1APIClient) setRequestUserAgent(r *http.Request) {
	r.Header.Set("User-Agent", c.userAgent)
}

func (c *ensapV1APIClient) Initialize() error {
	httpClient, err := c.newHTTPClient()
	if err != nil {
		return errors.Wrap(err, "initializing http client")
	}

	c.httpClient = httpClient
	return nil
}

func (c *ensapV1APIClient) buildFullEndpoint(endpoint string, versioned bool) string {
	s := fmt.Sprintf("https://%s%s", c.Endpoint, endpoint)
	if versioned {
		s = s + "/v1"
	}

	return s
}

func debugDumpHTTPResponse(funcName string, resp *http.Response, body bool) {
	dump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("fail to dump HTTP response for debug")
		return
	}

	log.WithFields(log.Fields{
		"response": fmt.Sprintf("%s", dump),
	}).Debugf("%s() HTTP response", funcName)
}

func debugDumpHTTPRequest(funcName string, req *http.Request, body bool) {
	dump, err := httputil.DumpRequestOut(req, body)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("fail to dump HTTP request for debug")
		return
	}

	log.WithFields(log.Fields{
		"request": fmt.Sprintf("%s", dump),
	}).Debugf("%s() HTTP request", funcName)
}
