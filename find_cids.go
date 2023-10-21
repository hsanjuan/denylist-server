package dls

import (
	"fmt"
	"regexp"

	"github.com/ipfs/go-cid"
	logging "github.com/ipfs/go-log/v2"
)

var logger = logging.Logger("dls")

var baseCharsets = map[string][]string{
	"base2":        {"0", "01"},
	"base8":        {"7", "01234567"},
	"base10":       {"9", "0123456789"},
	"base16":       {"f", "0123456789abcdef"},
	"base32hex":    {"v", "0123456789abcdefghijklmnopqrstuv"},
	"base32":       {"b", "abcdefghijklmnopqrstuvwxyz234567"},
	"base32z":      {"h", "ybndrfg8ejkmcpqxot1uwisza345h769"},
	"base36":       {"k", "0123456789abcdefghijklmnopqrstuvwxyz"},
	"base36upper":  {"K", "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"},
	"base58flickr": {"Z", "123456789abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ"},
	"base58btc":    {"z", "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"},
	"cidv0":        {"Qm", "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"},
	"base64":       {"m", "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"},
	"base64url":    {"u", "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"},
}

var cidRegexps = compileRegexps()

func compileRegexps() []*regexp.Regexp {
	var exps []*regexp.Regexp

	for base, charset := range baseCharsets {
		l := "6," // find multibase strings of at least 6 chars
		if base == "cidv0" {
			l = "44" // cidv0s have 44 char length, fixed.
		}
		regexpStr := fmt.Sprintf("%s[%s]{%s}", charset[0], charset[1], l)
		exps = append(exps, regexp.MustCompile(regexpStr))
	}

	return exps
}

// FindCIDs extracts CIDs from text. Returns a deduplicated list of CID
// strings as they are found in the text. The list is deduplicated based on
// their multihash, with later appereances of a CID triumphing over previous
// ones containing the same multihash.
func FindCIDs(s string) []string {
	cids := make(map[string]string)

	for _, exp := range cidRegexps {
		cidStrs := exp.FindAllString(s, -1)
		for _, cStr := range cidStrs {
			c, err := cid.Decode(cStr)
			if err != nil {
				logger.Warnf("could not decode %s as cid", cStr)
				continue
			}
			cids[string(c.Hash())] = cStr
		}
	}
	var results []string
	for _, v := range cids {
		results = append(results, v)
	}
	return results
}
