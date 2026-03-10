package harden

import "github.com/UT-CTF/landschaft/util"

func runHardenLockAccounts() {
	util.RunAndRedirectScript("harden/lock_accounts.ps1")
}
