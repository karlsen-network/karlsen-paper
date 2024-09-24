package paper

import (
	"github.com/karlsen-network/karlsend/v2/cmd/karlsenwallet/libkarlsenwallet"
	"github.com/karlsen-network/karlsend/v2/cmd/karlsenwallet/utils"
	"github.com/karlsen-network/karlsend/v2/domain/dagconfig"
	"github.com/karlsen-network/karlsen-paper/model"
	"github.com/pkg/errors"
	"github.com/tyler-smith/go-bip39"
	"fmt"
	"os"
	"bufio"
)

// Make sure we implement model.PaperWallet
var _ model.PaperAPI = &api{}

type api struct {
	dagParams *dagconfig.Params
}

func NewAPI(dagParams *dagconfig.Params) model.PaperAPI {
	return &api{dagParams: dagParams}
}

func (a *api) GenerateWallet() (model.PaperWallet, error) {

	fmt.Printf("Please enter your 24-word mnemonic or press enter for new:\n")
	reader := bufio.NewReader(os.Stdin)
	var mnemonics string
	var err error 

	mnemonics, err = utils.ReadLine(reader)

	if err != nil {
		return nil, err
	}

	if mnemonics == "" {
		fmt.Printf("Creating new mnemonics\n")
		mnemonics, err = libkarlsenwallet.CreateMnemonic()
	} else {
		if !bip39.IsMnemonicValid(string(mnemonics)) {
			return nil, errors.Errorf("mnemonic is invalid")
		}
	}

	return newWallet(a.dagParams, mnemonics)
}
