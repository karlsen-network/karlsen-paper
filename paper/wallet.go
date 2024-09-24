package paper

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/karlsen-network/karlsend/v2/cmd/karlsenwallet/keys"
	"github.com/karlsen-network/karlsend/v2/cmd/karlsenwallet/libkarlsenwallet"
	"github.com/karlsen-network/karlsend/v2/domain/dagconfig"
	"github.com/skip2/go-qrcode"
	"github.com/karlsen-network/karlsen-paper/model"
)

// Make sure we implement model.PaperWallet
var _ model.PaperWallet = &wallet{}

type wallet struct {
	dagParams *dagconfig.Params
	mnemonic  *model.MnemonicString
	keysFile  *keys.File
	password  string
}

func newWallet(dagParams *dagconfig.Params, mnemonic string) (model.PaperWallet, error) {
	mnemonicString := &model.MnemonicString{}
	copy(mnemonicString[:], strings.Split(mnemonic, " ")) // We assume it splits to 24 words

	password, err := generatePassword()
	if err != nil {
		return nil, err
	}

	keysFile, err := keys.NewFileFromMnemonic(dagParams, mnemonicString.String(), password)
	if err != nil {
		return nil, err
	}

	return &wallet{
		dagParams: dagParams,
		mnemonic:  mnemonicString,
		keysFile:  keysFile,
		password:  password,
	}, nil
}

func (w *wallet) Mnemonic() *model.MnemonicString {
	return w.mnemonic
}

func (w *wallet) Address(index int) (string, error) {
	path := fmt.Sprintf("m/%d/%d", libkarlsenwallet.ExternalKeychain, index)
	address, err := libkarlsenwallet.Address(w.dagParams, w.keysFile.ExtendedPublicKeys, 1, path, false)
	if err != nil {
		return "", err
	}
	return address.String(), nil
}

func (w *wallet) KPubKey() (string) {
	return w.keysFile.ExtendedPublicKeys[0]
}

func (w *wallet) AddressQR(index int) ([]byte, error) {
	address, err := w.Address(index)
	if err != nil {
		return nil, err
	}

	qr, err := qrcode.Encode(address, qrcode.High, 256)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return qr, nil
}

func (w *wallet) KPubKeyQR() ([]byte, error) {
	address := w.KPubKey()

	qr, err := qrcode.Encode(address, qrcode.High, 256)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return qr, nil
}
