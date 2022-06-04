# certissuer

This little tool can issue certificates from [Hashicorp Vault](https://vaultproject.io/). Originally, it has been written to guarantee certificate refresh (which didn't happen sometimes using automatic refresh capabilities of various Vault-aware tools), but it comes handy in a lot of situations.

## Usage

```shell
certissuer [global options]
```

where *global options* can be:

* --address <value> | -A <value>: Vault URL address (default: `$VAULT_ADDR`). Required.
* --token <value> | -T <value>: Vault token (default: `$VAULT_TOKEN`). Required.
* --cn <value> | -C <value>: common name used in issued certificate. Required.
* --pki <value> | -P <value>: [PKI mount](https://www.vaultproject.io/docs/secrets/pki#setup) in Vault. Required.
* --role <value> | -R <value>: [PKI role](https://www.vaultproject.io/docs/secrets/pki#configure-a-role) in Vault. Required.
* --ip <value> | -i <value>: add an IP-based Subject Alternative Name
* --alt <value> | -a <value>: add a DNS name-based Subject Alternative Name
* --ttl <value> | -t <value>: expiration duration. Sent to the server verbatim (default: 50h)
* --ssldir <value> | -d <value>: Directory to put certificates to (default: `/etc/ssl`)
* --no-verify | -k: Skips TLS verification while communicating to Vault (default: `$VAULT_SKIP_VERIFY`)
* --format <value> | -n <value>: Certificate format (see later). Default: "well-known"
* --help: shows help
* --version: shows version

## Formats

Certissuer can emit certificates in formats suit their purpose: providing certificates for HAProxy is different than doing it for NGINX. You can prefer storing all the certificates in a shared directory, or you might want to do it just how [Certbot](https://certbot.eff.org/) does.

Therefore, certissuer provides three formats:

* well-known: every cert goes into its own directory, and each data goes into their own well-known file name. Eg. CA goes into `ca.pem`, leaf certificate goes into `cert.pem`, certificate and intermediate certs go to `fullcert.pem`, and private key goes to `privkey.pem`. This is the default.
* named: this format saves full certificate chain (leaf and intermediate certificates concatenated) into `<CN>.pem`, and private key goes into `<CN>.key`.
* haproxy: this format generates a single file called `<CN>.pem`, containing leaf and intermediate certificates, and private key concatenated together.

## Any issues?

Open a ticket, perhaps a pull request. We support [GitHub Flow](https://guides.github.com/introduction/flow/). You might want to [fork](https://guides.github.com/activities/forking/) this project first.

## Code of Conduct

Small utilities are not commonly gather large communities, but we still require that discursion should be civil. Therefore, we adapt Contributor Covenant's Code of Conduct for this project.

## Licensing

SPDX-License-Identifier: BlueOak-1.0.0 OR MIT

This software is licensed under two licenses of your choice: Blue Oak Public License 1.0, or MIT Public License.
