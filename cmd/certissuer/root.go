package main

import (
	"strings"

	"github.com/julian7/certissuer/tlssave"
	"github.com/julian7/certissuer/vaultapi"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

type App struct {
	logger *zap.Logger
}

func NewApp(l *zap.Logger) *App {
	return &App{logger: l}
}

func (a *App) Command() *cli.App {
	return &cli.App{
		Name:            "certissuer",
		Usage:           "Issues TLS certificates from vault",
		ArgsUsage:       " ",
		HideHelpCommand: true,
		Version:         version,
		Action:          a.Action,
		Description: `This application creates new TLS certificates from  the vault server,
and stores it in well known places for your application to pick up.

By default, it saves generated certificates into well-known names, but it's
also possible to use provided Common Name as file names:

certificate     | well-known    | named         | haproxy
----------------+---------------+---------------+----------------
CA              | ca.pem        | (not written) | (not written)
certificate     | cert.pem      | (not written) | (not written)
cert+chain      | fullcert.pem  | <CN>.pem      | (not written)
private key     | privkey.pem   | <CN>.key      | (not written)
cert+chain+priv | (not written) | (not written) | <CN>.pem

To use either "named" or "haproxy" formats, provide "--format named" or
"--format haproxy" options.

Certificate chain contains all the certificates up to the root CA.`,

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "address",
				Aliases:  []string{"A"},
				EnvVars:  []string{"VAULT_ADDR"},
				Usage:    "Vault address",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "token",
				Aliases:  []string{"T"},
				EnvVars:  []string{"VAULT_TOKEN"},
				Usage:    "Vault token",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "cn",
				Aliases:  []string{"C"},
				Value:    "",
				Usage:    "Common Name (primary hostname)",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "pki",
				Aliases:  []string{"P"},
				Value:    "",
				Usage:    "PKI mount in Vault",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "role",
				Aliases:  []string{"R"},
				Value:    "",
				Usage:    "PKI role in Vault",
				Required: true,
			},
			&cli.StringSliceFlag{
				Name:    "ip",
				Aliases: []string{"i"},
				Value:   nil,
				Usage:   "IP SANs (Subject Alternative Names): bare IP addresses the cert should be valid",
			},
			&cli.StringSliceFlag{
				Name:    "alt",
				Aliases: []string{"a"},
				Value:   nil,
				Usage:   "SANs (Subject Alternative Names): other names the cert should be valid",
			},
			&cli.StringFlag{
				Name:    "ttl",
				Aliases: []string{"t"},
				Value:   "50h",
				Usage:   "Time-to live, eg. expiration duration. Sent to vault server verbatim.",
			},
			&cli.StringFlag{
				Name:    "ssldir",
				Aliases: []string{"d"},
				Value:   "/etc/ssl",
				Usage:   "Directory to put certificates to",
			},
			&cli.BoolFlag{
				Name:    "no-verify",
				Aliases: []string{"k"},
				EnvVars: []string{"VAULT_SKIP_VERIFY"},
				Value:   false,
				Usage:   "Skip TLS verification while communicating to Vault. Required to rearm vault.",
			},
			&cli.StringFlag{
				Name:    "format",
				Aliases: []string{"n"},
				Value:   "well-known",
				Usage:   "Certificate format. See long help for details.",
			},
		},
	}
}

func (a *App) Action(ctx *cli.Context) error {
	pki := ctx.String("pki")
	role := ctx.String("role")
	cn := ctx.String("cn")

	ipsans := ctx.StringSlice("ip")
	alts := ctx.StringSlice("alt")
	ttl := ctx.String("ttl")
	ssldir := ctx.String("ssldir")
	format := ctx.String("format")

	formatType := 0
	switch format {
	case "well-known":
		formatType = tlssave.PersistFunctional
	case "named":
		formatType = tlssave.PersistNamed
	case "haproxy":
		formatType = tlssave.PersistHAProxy
	default:
		a.logger.Warn("invalid format type", zap.String("format", format))
		return nil
	}

	api, err := vaultapi.New(ctx.String("address"), ctx.String("token"), !ctx.Bool("no-verify"))

	if err != nil {
		return err
	}

	// step 1: check vault
	err = api.GetHealth()
	if err != nil {
		return err
	}
	// step 2: pull cert
	req := &vaultapi.CertRequest{
		CommonName: cn,
		AltNames:   strings.Join(alts, ","),
		IPSans:     strings.Join(ipsans, ","),
		TTL:        ttl,
	}
	resp, err := api.IssueCert(pki, role, req)
	if err != nil {
		return err
	}
	if len(resp.Warnings) > 0 {
		a.logger.Warn(resp.Warnings)
	}
	// step 3: save files
	persister, err := tlssave.New(
		formatType,
		cn,
		resp.Data.Certificate,
		resp.Data.IssuingCA,
		resp.Data.CAChain,
		resp.Data.PrivateKey,
	)
	if err != nil {
		return err
	}
	if err := persister.Persist(ssldir); err != nil {
		return err
	}

	return nil
}
