/*
   Debian cli eXtension for Artisan Enterprise
   © 2022 SouthWinds Tech Ltd

   Original Source taken from Flufik by Eduard Gevorkyan
   https://github.com/egevorkyan/flufik/blob/master/LICENSE
*/

package apt

import (
	"bytes"
	"crypto"
	"fmt"
	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/clearsign"
	"github.com/ProtonMail/go-crypto/openpgp/packet"
	"io"
	"os"
	"path/filepath"
)

func readPrivateKey(key, passcode []byte) (entity *openpgp.Entity, err error) {
	var entityList openpgp.EntityList
	entityList, err = openpgp.ReadArmoredKeyRing(bytes.NewReader(key))
	if err != nil {
		return nil, fmt.Errorf("decoding armored PGP keyring failure %w", err)
	}
	entity = entityList[0]
	if entity.PrivateKey == nil {
		return nil, fmt.Errorf("no private key")
	}
	if entity.PrivateKey.Encrypted {
		if string(passcode) == "" {
			return nil, fmt.Errorf("key encrypted, passphrase not provided")
		}
		if err = entity.PrivateKey.Decrypt(passcode); err != nil {
			return nil, fmt.Errorf("failure decrypting private key: %w", err)
		}
		for _, subKey := range entity.Subkeys {
			if subKey.PrivateKey != nil {
				if err = subKey.PrivateKey.Decrypt(passcode); err != nil {
					return nil, fmt.Errorf("failure decrypting sub private key: %w", err)
				}
			}
		}
	}
	return entity, nil
}

func signRelease(signingKey, passcode []byte, fileName string) error {
	key, err := readPrivateKey(signingKey, passcode)
	if err != nil {
		return err
	}
	workingDirectory := filepath.Dir(fileName)
	releaseFile, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("error opening release file (%s) for writing: %s", fileName, err)
	}
	releaseGpg, err := os.Create(filepath.Join(workingDirectory, "Release.gpg"))
	if err != nil {
		return fmt.Errorf("error creating Release.pgp file for writing: %s", err)
	}
	defer releaseGpg.Close()
	if err = openpgp.ArmoredDetachSign(releaseGpg, key, releaseFile, &packet.Config{
		DefaultHash: crypto.SHA256,
	}); err != nil {
		return fmt.Errorf("armored detached sign failure: %s", err)
	}
	releaseFile.Seek(0, 0)
	var inlineRelease *os.File
	inlineRelease, err = os.Create(filepath.Join(workingDirectory, "InRelease"))
	if err != nil {
		return fmt.Errorf("error creating InRelease file for writing: %s", err)
	}
	defer inlineRelease.Close()
	var writer io.WriteCloser
	writer, err = clearsign.Encode(inlineRelease, key.PrivateKey, nil)
	if err != nil {
		return fmt.Errorf("error signing InRelease file : %s", err)
	}
	if _, err = io.Copy(writer, releaseFile); err != nil {
		return err
	}
	return writer.Close()
}
