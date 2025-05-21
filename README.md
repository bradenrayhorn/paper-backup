# paper-backup

Backup your files and secrets to paper.

paper-backup is distributed as a single, self-contained HTML file. It can be used offline for maximum security.

Demo available at https://paperbackup.hydrated.dev.

There is a minimal, but functional UI.

## Usage

TBD

## Features

### File backup

Upload a file and enter a passphrase. paper-backup will compress and encrypt your file,
then generate a printable QR code containing your file. The security of the encryption will
of course depend on the length of your passphrase.

To recover the original file, scan the QR code using paper-backup and enter the passphrase. The
original file will be available again.

Some example uses:

- 2FA backup codes
- Cryptocurrency wallet
- GPG or SSH keys

**Note**: Due to the size limitations of QR codes, this only works for a few KB of data.

### Shamir split secret

Uses [Shamir's secret sharing algorithm](https://en.wikipedia.org/wiki/Shamir%27s_secret_sharing) to split a secret
into multiple parts. A configurable number of parts is required to reassemble the original secret.

For example, if you split a secret into 5 parts with a minimum threshold of 3 parts required to recover,
any 3 of the original 5 parts can be put together to recover the original secret. The secret parts can be
distributed throughout your local area at a friend's house, work, parent's house, etc.

Enter a small secret (a paragraph or two of text), choose your number of parts and threshold, and paper-backup
will output the requested number of secret parts. Each part is intended to be printed double-sided,
one side containst the secret in the form of a QR code, the other as text (in case of QR code damage).

To recover, scan or type the required number of parts into paper-backup. The original secret will be
displayed.

Some example uses:

- Important passwords
- Encryption keys

This method of backup provides redundancy. A few parts of the secret could be destroyed in a natural disaster,
but the secret is still recoverable with the remaining parts (if your parts were distributed across a large geographic area).

At the same time, this backup provides some additional security over writing passwords on paper. An adversary would
need to acquire a number of paper parts to reassemble the secret, rather than just one.

## Credits

Uses Hashicorp Vault's shamir secret sharing [implementation](https://github.com/hashicorp/vault/tree/e31d45514d4314fdb90f169a66523c4f8feb789c/shamir).
Code used is in the `shamir/` package.

The [zxing-cpp](https://github.com/zxing-cpp/zxing-cpp) library is compiled to Web Assembly and used for QR code generation.

## Security

This project has not been audited by a third party.

Anyone using paper-backup should consider their individual threat model and use at their own risk.

The code is fairly straightforward if you'd like to read through it. All cryptography is performed
using the Go standard library (and official Go sub-repositories).

