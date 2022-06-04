package tlssave

import (
	"fmt"
	"path"
)

type TLSNamed struct {
	*bundle
}

func (t *TLSNamed) Persist(ssldir string) error {
	if err := checkDir(ssldir); err != nil {
		return fmt.Errorf("SSL dir: %w", err)
	}
	cn := t.bundle.CommonName
	err := t.bundle.saveFullCert(path.Join(ssldir, fmt.Sprintf("%s.pem", cn)))
	if err != nil {
		return err
	}
	return t.bundle.savePrivKey(path.Join(ssldir, fmt.Sprintf("%s.key", cn)))
}
