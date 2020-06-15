package cli

import (
	"fmt"
	"github.com/Songmu/prompter"
	"github.com/spf13/cobra"
	"github.com/sudachen/smwlt/fu"
	api "github.com/sudachen/smwlt/node/api.v1"
	"github.com/sudachen/smwlt/verbose"
	"github.com/sudachen/smwlt/wallet"
	"github.com/sudachen/smwlt/wallet/legacy"
	"github.com/sudachen/smwlt/wallet/modern"
	"os"
	"path/filepath"
	"strings"
)

const MajorVersion = 1
const MinorVersion = 0
const keyValueFormat = "%-8s %v\n"

var mainCmd = &cobra.Command{
	Use:           "smwlt",
	Short:         fmt.Sprintf("Spacemesh CLI Wallet %v.%v (https://github.com/sudachen/smwlt)", MajorVersion, MinorVersion),
	SilenceErrors: true,
}

var optWalletFile = mainCmd.PersistentFlags().StringP("wallet-file", "f", "", "use wallet filename")
var optWalletName = mainCmd.PersistentFlags().StringP("wallet-name", "n", "", "select wallet by name")
var optWalletDir = mainCmd.PersistentFlags().StringP("wallet-dir", "d", modern.DefaultDirectory(), "use wallet dir")
var optLegacy = mainCmd.PersistentFlags().BoolP("legacy", "l", false, "use legacy unencrypted file format")
var optPassword = mainCmd.PersistentFlags().StringP("password", "p", "", "wallet unlock password")
var optEndpoint = mainCmd.PersistentFlags().StringP("endpoint", "e", api.DefaultEndpoint, "host:port to connect mesh node")
var optYes = mainCmd.PersistentFlags().BoolP("yes", "y", false, "auto confirm")
var OptTrace = mainCmd.PersistentFlags().BoolP("trace", "x", false, "backtrace on panic")

func Main() {
	verbose.VerboseOptP = mainCmd.PersistentFlags().BoolP("verbose", "v", false, "be verbose")
	mainCmd.AddCommand(
		cmdInfo,
		cmdSend,
		cmdTxs,
		cmdNet,
		cmdHexSign,
		cmdTextSign,
		cmdCoinbase,
		cmdNew,
	)

	if err := mainCmd.Execute(); err != nil {
		panic(fu.Panic(err, 1))
	}
}

func unlock(w wallet.Wallet, passw *[]string, interactive bool) bool {
	for _, p := range *passw {
		if e := w.Unlock(p); e == nil {
			return true
		}
	}
	if interactive {
		fmt.Printf("Unlocking wallet %v\n", w.DisplayName())
		p := prompter.Password("Enter password [leave empty to skip]")
		if p != "" {
			if e := w.Unlock(p); e == nil {
				fmt.Println("Wallet unlocked")
				*passw = append(*passw, p)
				return true
			} else {
				fmt.Println("Wrong password!")
			}
		} else {
			fmt.Println("Wallet skipped")
		}
	}
	return false
}

func loadWallet() (w []wallet.Wallet) {
	if *optLegacy {
		w = []wallet.Wallet{legacy.Wallet{Path: *optWalletFile}.LuckyLoad()}
		// unencrypted
	} else {
		w = []wallet.Wallet{}
		wx := []wallet.Wallet{}
		passw := []string{}
		if *optPassword != "" {
			passw = append(passw, *optPassword)
		}
		if *optWalletFile != "" {
			wx = []wallet.Wallet{modern.Wallet{Path: *optWalletFile}.LuckyLoad()}
		} else {
			if err := filepath.Walk(*optWalletDir, func(path string, info os.FileInfo, err error) error {
				base := filepath.Base(path)
				if strings.HasPrefix(base, "my_wallet_") {
					verbose.Printfln("opening wallet file '%v'", base)
					wal, err := modern.Wallet{Path: path}.Load()
					if err == nil {
						if *optWalletName == "" ||
							strings.HasPrefix(strings.ToLower(wal.DisplayName()), strings.ToLower(*optWalletName)) {
							wx = append(wx, wal)
						}
					} else {
						verbose.Printfln("failed to open with error: %v", err.Error())
					}
				}
				return nil
			}); err != nil {
				panic(fu.Panic(err))
			}
		}
		for _, x := range wx {
			if unlock(x, &passw, *optPassword == "") {
				w = append(w, x)
			}
		}
		if len(w) == 0 && *optPassword != "" {
			panic(fu.Panic(fmt.Errorf("there is nothing to unlock, wrong password(?)")))
		}
	}
	return
}

type Client struct {
	*api.ClientAgent
	DefaultGasLimit uint64
	DefaultFee      uint64
}

func newClient() Client {
	return Client{
		ClientAgent:     api.Client{Endpoint: *optEndpoint, Verbose: verbose.Printfln}.New(),
		DefaultGasLimit: api.DefaultGasLimit,
		DefaultFee:      api.DefaultFee,
	}
}