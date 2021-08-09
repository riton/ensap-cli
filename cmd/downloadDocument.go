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
package cmd

import (
	"io"
	"os"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var downloadDocumentCmdFlags = struct {
	DownloadDestinationFile string
}{}

// downloadDocumentCmd represents the downloadDocument command
var downloadDocumentCmd = &cobra.Command{
	Use:   "download-remuneration-paie",
	Short: "Download a remuneration-paie document ID",
	Long:  ``,
	RunE:  downloadDocumentCmdRunE,
	Args:  cobra.ExactArgs(1),
}

func init() {
	rootCmd.AddCommand(downloadDocumentCmd)

	downloadDocumentCmd.Flags().StringVarP(&downloadDocumentCmdFlags.DownloadDestinationFile, "out", "o", "", "destination file")
}

func downloadDocumentCmdRunE(cmd *cobra.Command, args []string) error {
	apiClient := cmdNewEnsapAPIClient()
	if err := apiClient.Initialize(); err != nil {
		return err
	}

	if err := apiClient.Login(); err != nil {
		return errors.Wrap(err, "login failed")
	}

	defer func() {
		if err := apiClient.Logout(); err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("fail to logout")
		}
	}()

	document, err := apiClient.DownloadRemunerationPaie(args[0])
	if err != nil {
		return errors.Wrap(err, "downloading remuneration paie document")
	}

	var destWriter io.WriteCloser
	if downloadDocumentCmdFlags.DownloadDestinationFile == "-" {
		destWriter = os.Stdout
	} else {
		documentFilename := document.Filename
		if downloadDocumentCmdFlags.DownloadDestinationFile != "" {
			documentFilename = downloadDocumentCmdFlags.DownloadDestinationFile
		}
		fd, err := os.Create(documentFilename)
		if err != nil {
			return errors.Wrapf(err, "creating output file %q", document.Filename)
		}

		destWriter = fd
	}

	defer func() {
		document.ReadCloser.Close()
		destWriter.Close()
	}()

	if _, err := io.Copy(destWriter, document.ReadCloser); err != nil {
		return errors.Wrap(err, "writing file")
	}

	return nil
}
