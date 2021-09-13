package objects

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
	"time"

	"github.com/pkg/errors"
)

// RemunerationPaieDocument is a document as returned by
// the https://ensap.gouv.fr/prive/remunerationpaie/v1 endpoint
type RemunerationPaieDocument struct {
	DocumentUUID string       `json:"documentUuid"`
	Libelle1     string       `json:"libelle1"`
	Libelle2     string       `json:"libelle2"`
	Libelle3     string       `json:"libelle3"`
	NomDocument  string       `json:"nomDocument"`
	DateDocument documentTime `json:"dateDocument"`
	Annee        int          `json:"annee"`
	Icone        string       `json:"icone"`
	LibelleIcone string       `json:"libelleIcone"`
}

type documentTime struct {
	time.Time
}

// UnmarshalJSON is required to avoid
// decoding JSON response: parsing time ""2021-07-01T00:00:00.000+0200"" as ""2006-01-02T15:04:05Z07:00"": cannot parse "+0200"" as "Z07:00"
func (d *documentTime) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	t, err := time.Parse("2006-01-02T15:04:05.000Z0700", s)
	if err != nil {
		return errors.Wrap(err, "parsing date")
	}

	*d = documentTime{t}

	return nil
}

func (d documentTime) String() string {
	return d.Time.Format(time.RFC822)
}
