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
	"html/template"
	"io"
	"os"

	"github.com/Masterminds/sprig"
	"github.com/pkg/errors"
	"github.com/riton/ensap-cli/api/v1/objects"
	"github.com/riton/ensap-cli/timeutils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listDocumentsCmd represents the listDocuments command
var listDocumentsCmd = &cobra.Command{
	Use:           "list-remuneration-paie",
	Short:         "List all remuneration-paie documents",
	Long:          ``,
	RunE:          listDocumentsCmdRunE,
	SilenceUsage:  true,
	SilenceErrors: true,
}

var listDocumentsCmdFlags = struct {
	Year int
}{}

func init() {
	rootCmd.AddCommand(listDocumentsCmd)
	listDocumentsCmd.Flags().IntVarP(&listDocumentsCmdFlags.Year, "year", "y", timeutils.GetCurrentYear(), "year to list documents for")
}

func listDocumentsCmdRunE(cmd *cobra.Command, args []string) error {
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

	documents, err := apiClient.ListRemunerationPaieDocumentsByYear(listDocumentsCmdFlags.Year)
	if err != nil {
		return err
	}

	// user can specify its own template
	renderingTemplate := viper.GetString("rendering.templates.list-remuneration-paie-documents")
	if renderingTemplate == "" {
		renderingTemplate = listRemunerationPaieDocumentsDefaultTpl
	}

	return renderDocumentsTemplateToWriter(os.Stdout, renderingTemplate, documents)
}

const (
	listRemunerationPaieDocumentsDefaultTpl = `Remuneration paie documents:
{{- range $index, $document := . }}
{{- with $document }}
* {{ .NomDocument }}:
  - Document UUID: {{ .DocumentUUID }}
  - Document date: {{ .DateDocument }}
  - Labels:
    - {{ .Libelle1 }}
    - {{ .Libelle2 }}
    - {{ .Libelle3 }}
{{- end }}
{{- end }}
`
)

func renderDocumentsTemplateToWriter(w io.Writer, tpl string, documents []objects.RemunerationPaieDocument) error {
	funcMap := sprig.FuncMap()

	customFuncMap := template.FuncMap{}

	t := template.Must(template.New("t1").Funcs(funcMap).Funcs(customFuncMap).Parse(tpl))
	if err := t.Execute(w, documents); err != nil {
		return errors.Wrap(err, "executing template")
	}
	return nil
}
