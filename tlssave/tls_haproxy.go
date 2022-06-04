package tlssave

import (
	"fmt"
	"path"
)

type TLSHAProxy struct {
	*bundle
}

func (t *TLSHAProxy) Persist(ssldir string) error {
	if err := checkDir(ssldir); err != nil {
		return fmt.Errorf("SSL dir: %w", err)
	}
	return t.bundle.saveFile(
		path.Join(ssldir, fmt.Sprintf("%s.pem", t.bundle.CommonName)),
		fmt.Sprintf(
			"%s\n%s\n%s\n",
			t.bundle.Cert,
			t.bundle.Chain,
			t.bundle.PrivKey,
		),
		0600,
	)
}
