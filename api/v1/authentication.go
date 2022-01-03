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
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
	"github.com/riton/ensap-cli/api"
	log "github.com/sirupsen/logrus"
)

const (
	// AuthentificationEndpointAuthOKCode is the code returned by the "authentification"
	// endpoint as in { "code":60, "message":"Authentification OK"}
	AuthentificationEndpointAuthOKCode = 60
)

type authentificationEndpointReponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (c *ensapV1APIClient) Login() error {
	formData := url.Values{}
	formData.Set("identifiant", c.username)
	formData.Set("secret", c.password)

	endpoint := c.buildFullEndpoint(api.AuthentificationEndpoint, false)

	log.Debugf("logging in using endpoint %s", endpoint)

	req, err := http.NewRequestWithContext(c.ctx, http.MethodPost, endpoint, strings.NewReader(formData.Encode()))
	if err != nil {
		return errors.Wrap(err, "creating request object")
	}
	c.setRequestUserAgent(req)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	// if c.cfg.Debug {
	// 	debugDumpHTTPRequest("Login", req, true)
	// }

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return errors.Wrapf(err, "performing HTTP request to %s", endpoint)
	}
	defer resp.Body.Close()

	if c.cfg.Debug {
		debugDumpHTTPResponse("Login", resp, true)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad response code %d from endpoint %s", resp.StatusCode, endpoint)
	}

	var authStatus authentificationEndpointReponse
	if err := json.NewDecoder(resp.Body).Decode(&authStatus); err != nil {
		return errors.Wrap(err, "decoding JSON response")
	}

	if authStatus.Code != AuthentificationEndpointAuthOKCode {
		return fmt.Errorf("bad API response code %d with message %q. %d was expected", authStatus.Code, authStatus.Message, AuthentificationEndpointAuthOKCode)
	}

	log.Debug("successfully authenticated")

	return nil
}

func (c *ensapV1APIClient) Logout() error {
	endpoint := c.buildFullEndpoint(api.DeconnexionEndpoint, true)

	req, err := http.NewRequestWithContext(c.ctx, http.MethodPost, endpoint, nil)
	if err != nil {
		return errors.Wrap(err, "creating request object")
	}
	c.setRequestUserAgent(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return errors.Wrapf(err, "performing HTTP request to %s", endpoint)
	}
	defer resp.Body.Close()

	if c.cfg.Debug {
		debugDumpHTTPResponse("Logout", resp, true)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad response code %d from endpoint %s", resp.StatusCode, endpoint)
	}

	log.Debug("successfully disconnected")

	return nil
}
