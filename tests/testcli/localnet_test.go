// +build !windows

package testcli

import (
	"fmt"
	"github.com/sudachen/smwlt/tests/expect"
	"gotest.tools/assert"
	"io/ioutil"
	"strings"
	"testing"
)

const endpoint = "localhost:19090"
const tapAddress = "0x6Ef675c1BB32E9021844193A163FaC231513149A"
const tapBalance = 10000

const accountsJson = `{"almog":{"pubkey":"0xb8110cfeB1f01011E118BdB93F1Bb14D2052c276","privkey":"48dc32104006e7a85c53eeb80cd54dd218d8b66908ad5a9a12bee7b6c01461904d05cfede9928bcd225c008db8110cfeb1f01011e118bdb93f1bb14d2052c276"},"anton":{"pubkey":"0x66Bf7ef2D1Da7F0b391D1F1364f1D695929df617","privkey":"af94f1477c0325e5b790178e110cb1f51537589f1726edaaf8cacd385145d289db58184012f26c405bff2d8866bf7ef2d1da7f0b391d1f1364f1d695929df617"},"barak":{"pubkey":"0x887A595e41B097AcD0A75d65Ed8b8C6Fa739D297","privkey":"ded99e45c7869465bf0690ebc39d700785147cf99d0e8c7507f839053519b73b097598942e44919cf7d11499887a595e41b097acd0a75d65ed8b8c6fa739d297"},"gavrad":{"pubkey":"0x437E320d792772AbA8F459f80E18a45ae754112d","privkey":"19bf601036b2a1770ff1ce359d76119a03fbfee91d9fc2d870bc37e34dc2d9890dc90fe42d96e302ae122aa3437e320d792772aba8f459f80e18a45ae754112d"},"tap":{"pubkey":"0x6Ef675c1BB32E9021844193A163FaC231513149A","privkey":"a3db6583a3e989525fac730814f85d9c9fef3c00f4ed1c4c56c33ac560c94950891da146767aa80e3ce3ef826ef675c1bb32e9021844193a163fac231513149a"},"test":{"pubkey":"0xE9a546c21E9E0817C0fDc2E8Daa2668d926bAeCA","privkey":"b21489741cb1a168fda9fcad25fb39218bfbe8f80c85ce0f1e65777f78b3fa9789d491ae29a13be054de739ae9a546c21e9e0817c0fdc2e8daa2668d926baeca"},"test2":{"pubkey":"0x206E49951487c0E3004e82c34f5746E90FDbf9AE","privkey":"3bd3eca43eb67f40dc1ca52a039f7a1fc100d7275ca4425fefae39cb3a1fa75d1e35150f2f1d372794d68970206e49951487c0e3004e82c34f5746e90fdbf9ae"},"test3":{"pubkey":"0xBAF58761Adf1416DB5679B097C60A7E86E4720CD","privkey":"75ccc3a74edaa57457c259cb25e83116cc09127cc44806938b481bc617bbd23efd87cee3cc8ad387cba6dbe8baf58761adf1416db5679b097c60a7e86e4720cd"},"yosher":{"pubkey":"0xF94Abe7Ba1428df096e13e903Ef5B9dF85d520e1","privkey":"8bab7df1f55b1f0160e130abc1e3d7c5985345a3c1f8e6d3f7b4b1defad3b67439a27e846f7e9783cd8fcae0f94abe7ba1428df096e13e903ef5b9df85d520e1"}}`
const legacyWallet = "test-accounts.json"

func init() {
	if err := ioutil.WriteFile(legacyWallet, []byte(accountsJson), 0666); err != nil {
		panic(err)
	}
}


func Test_CmdNet1(t *testing.T) {
	testCLI(t, func(t *testing.T, pty *expect.Pty) {
		pty.Host.SkipToExpect(`Node status:`)
		pty.Host.Expect(`Synced:\s*(false|true)`)
		pty.Host.Expect(`Synced layer:\s*\d+`)
		pty.Host.Expect(`Current layer:\s*\d+`)
		pty.Host.Expect(`Peers:\s*\d+`)
		pty.Host.Expect(`Min peers:\s*\d+`)
		pty.Host.Expect(`Max peers:\s*\d+`)
		pty.Host.Expect(`Data directory:\s*(\w|/|-|.)+`)
		pty.Host.Expect(`Mining status:\s*\w+`)
		pty.Host.Expect(`Coinbase:\s*0x[A-Fa-f0-9]+`)
	}, "net", "-e", endpoint)
}

func Test_CmdInfo1(t *testing.T) {
	// tap does not mine on the local testnet
	testCLI(t, func(t *testing.T, pty *expect.Pty) {
		pty.Host.SkipToExpect(`Account \w+ \[legacyWallet\(`+legacyWallet+`\)\]:`)
		r := pty.Host.ExpectGet(`Address:\s*(0x[A-Fa-f0-9]+)`)
		assert.Equal(t, strings.ToLower(r[0]), strings.ToLower(tapAddress))
		pty.Host.Expect(`Balance:\s*`+fmt.Sprintf("%d",tapBalance))
	}, "info", "-l", "-f", legacyWallet, "-e", endpoint, "tap")
}

