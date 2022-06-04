package tlssave

import (
	"fmt"
	"os"
	"strings"
)

const (
	_ = iota
	PersistFunctional
	PersistNamed
	PersistHAProxy
)

type TLSPersister interface {
	Persist(ssldir string) error
}

//bundle contains all the data a TLS certificate should contain
type bundle struct {
	CommonName string
	Cert       string
	CA         string
	Chain      []string
	PrivKey    string
}

//New creates a new bundle
func New(persistType int, cn, cert, ca string, cachain []string, privkey string) (TLSPersister, error) {
	bundle := &bundle{
		CommonName: cn,
		Cert:       cert,
		CA:         ca,
		Chain:      cachain,
		PrivKey:    privkey,
	}
	switch persistType {
	case PersistFunctional:
		return &TLSFunctional{bundle: bundle}, nil
	case PersistNamed:
		return &TLSNamed{bundle: bundle}, nil
	case PersistHAProxy:
		return &TLSHAProxy{bundle: bundle}, nil
	}
	return nil, fmt.Errorf("unknown persisting type %d", persistType)
}

func checkDir(dir string) error {
	st, err := os.Stat(dir)
	if err != nil {
		return err
	}
	if !st.Mode().IsDir() {
		return fmt.Errorf("%s is not a directory", dir)
	}
	return nil
}

func (bundle *bundle) saveCA(fn string) error {
	return bundle.saveFile(fn, fmt.Sprintf("%s\n", bundle.CA), 0644)
}
func (bundle *bundle) saveCert(fn string) error {
	return bundle.saveFile(fn, fmt.Sprintf("%s\n", bundle.Cert), 0644)
}
func (bundle *bundle) saveFullCert(fn string) error {
	return bundle.saveFile(
		fn,
		fmt.Sprintf(
			"%s\n%s\n",
			bundle.Cert,
			strings.Join(bundle.Chain, "\n"),
		),
		0644,
	)
}
func (bundle *bundle) savePrivKey(fn string) error {
	return bundle.saveFile(fn, fmt.Sprintf("%s\n", bundle.PrivKey), 0600)
}

func (bundle *bundle) saveFile(fn, data string, perms os.FileMode) error {
	fd, err := os.OpenFile(fn, os.O_WRONLY|os.O_CREATE, perms)
	if err != nil {
		return fmt.Errorf("cannot open %s for writing: %w", fn, err)
	}
	defer fd.Close()
	_, err = fd.WriteString(data)
	if err != nil {
		return fmt.Errorf("cannot write into %s: %w", fn, err)
	}
	return nil
}
