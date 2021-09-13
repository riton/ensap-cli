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
	"mime"
	"net/http"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/riton/ensap-cli/api"
	"github.com/riton/ensap-cli/api/v1/objects"
)

func (c *ensapV1APIClient) ListRemunerationPaieDocumentsByYear(year int) ([]objects.RemunerationPaieDocument, error) {
	var documents []objects.RemunerationPaieDocument
	endpoint := c.buildFullEndpoint(api.RemunerationPaieEndpoint, true)

	req, err := http.NewRequestWithContext(c.ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return documents, errors.Wrap(err, "creating request object")
	}

	c.setRequestUserAgent(req)

	query := req.URL.Query()
	query.Add("annee", fmt.Sprintf("%d", year))
	req.URL.RawQuery = query.Encode()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return documents, errors.Wrapf(err, "performing HTTP request to %s", endpoint)
	}
	defer resp.Body.Close()

	if c.cfg.Debug {
		debugDumpHTTPRequest("ListRemunerationPaieDocumentsByYear", req, true)
		debugDumpHTTPResponse("ListRemunerationPaieDocumentsByYear", resp, true)
	}

	if resp.StatusCode != http.StatusOK {
		return documents, fmt.Errorf("bad response code %d from endpoint %s", resp.StatusCode, endpoint)
	}

	if err := json.NewDecoder(resp.Body).Decode(&documents); err != nil {
		return documents, errors.Wrap(err, "decoding JSON response")
	}

	return documents, nil
}

const (
	applicationPDFContentType = "application/pdf"
)

func (c *ensapV1APIClient) DownloadRemunerationPaie(documentUUID string) (objects.DownloadRemunerationPaieDocument, error) {
	var document objects.DownloadRemunerationPaieDocument

	endpoint := c.buildFullEndpoint(api.TelechargerRemunerationPaieEndpoint, true)

	req, err := http.NewRequestWithContext(c.ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return document, errors.Wrap(err, "creating request object")
	}

	c.setRequestUserAgent(req)

	query := req.URL.Query()
	query.Add("documentUuid", documentUUID)
	req.URL.RawQuery = query.Encode()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return document, errors.Wrapf(err, "performing HTTP request to %s", endpoint)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return document, fmt.Errorf("bad HTTP status code %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType != applicationPDFContentType {
		resp.Body.Close()
		return document, fmt.Errorf("bad content type %s, was expecting %s", contentType, applicationPDFContentType)
	}

	filename, err := httpRespGetAttachmentFilename(resp)
	if err != nil {
		resp.Body.Close()
		return document, errors.Wrap(err, "getting attachment filename")
	}

	if err := validateAttachmentFilename(filename, []string{".pdf"}); err != nil {
		resp.Body.Close()
		return document, errors.Wrapf(err, "invalid document filename %q", filename)
	}

	document.Filename = filename
	document.ReadCloser = resp.Body

	return document, nil
}

const (
	validateAttachmentFilenameMaxLen = 64
)

func validateAttachmentFilename(filename string, allowedExtensions []string) error {
	if len(filename) > validateAttachmentFilenameMaxLen {
		return fmt.Errorf("attachment filename %q is too long (max allowed is %d", filename, validateAttachmentFilenameMaxLen)
	}

	extension := filepath.Ext(filename)
	for _, allowedExtension := range allowedExtensions {
		if extension == allowedExtension {
			return nil
		}
	}

	return fmt.Errorf("file extension %s is not allowed", extension)
}

func httpRespGetAttachmentFilename(resp *http.Response) (string, error) {
	cd := resp.Header.Get("Content-Disposition")

	disposition, params, err := mime.ParseMediaType(cd)
	if err != nil {
		return "", errors.Wrap(err, "parsing media type")
	}

	if disposition != "attachment" {
		return "", fmt.Errorf("bad mediaType %s. 'attachment' was expected", disposition)
	}

	filename, ok := params["filename"]
	if !ok {
		return "", fmt.Errorf("filename missing from attachment")
	}

	return filename, nil
}
