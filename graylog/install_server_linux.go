package graylog

import (
	"path/filepath"

	"github.com/UT-CTF/landschaft/util"
)

func installServer(tlsPublicChain string, tlsPrivateKey string) {
	absPublicChain, err := filepath.Abs(tlsPublicChain)
	if err != nil {
		panic(err)
	}
	absPrivateKey, err := filepath.Abs(tlsPrivateKey)
	if err != nil {
		panic(err)
	}

	util.RunAndRedirectScript("graylog/install_server.sh", absPublicChain, absPrivateKey)
}
