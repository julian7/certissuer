package tlssave

import (
	"fmt"
	"path"
)

type TLSFunctional struct {
	*bundle
}

func (t *TLSFunctional) Persist(ssldir string) error {
	if err := checkDir(ssldir); err != nil {
		return fmt.Errorf("SSL dir: %w", err)
	}
	err := t.bundle.saveCA(path.Join(ssldir, "ca.pem"))
	if err != nil {
		return err
	}
	err = t.bundle.saveCert(path.Join(ssldir, "cert.pem"))
	if err != nil {
		return err
	}
	err = t.bundle.saveFullCert(path.Join(ssldir, "fullcert.pem"))
	if err != nil {
		return err
	}
	return t.bundle.savePrivKey(path.Join(ssldir, "privkey.pem"))
}
