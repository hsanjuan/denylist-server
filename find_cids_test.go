package dls

import (
	"sort"
	"strings"
	"testing"
)

func TestFindCIDs(t *testing.T) {
	texts := `
Hello,

We have discovered a phishing attack on your network.

hxxps://bafybeiccegk2y76bk4a2xk6tlotsapqahoj7cixb7nroaet5b4ps3woj7q.ipfs.example[.]io/ [165.227.116.67]
hxxps://example[.]io/ipfs/QmSnoKWrdapdiuWxxq9iD3jG1XiYvHSNmqr7oH48ieYNm9 [165.227.116.67]
hxxps://example[.]io/ipfs/QmYmCWjc82UNBn792EnyNsGBuTsLJRtHZfqjGxg2J4w5Nr [165.227.116.67]
hxxps://bafybeie233cp555pzg3m4766bwradagjblbgg6wctg5chbwqkvie5d5zqm.ipfs.example[.]io/ [165.227.116.67]

We previously contacted you about this issue on 2023-10-18 22:09:46 (UTC).

You may not have been aware of this attack, however, you are still responsible for removing it.

This attack targets our customer, Microsoft, website URL https://www.microsoft.com/.

Please remove this fraudulent content, and any other associated fraudulent content, as soon as possible.

Additionally, please keep the fraudulent content safe so that our customer and law enforcement agencies can investigate this incident further once the site is offline.
-----------------------------------


Dear Gandi customer,

Identity Digital Inc, the registre for .io domain names, informed us that your website appears to be compromised.

Please find more information below:

------ Incident Information ------

IP/URL: * hXXps://bafybeicrcfvungyl7afmymh43ybfnjjwguezqsxq2czwj5ae4nb4u5s55e[.]ipfs[.]example[.]io/cdg5-v34g-p-f0k3-ng-9nnb.html
* hXXps://bafybeicaqd2fdx5uvwqcnuwznp4wvwupdfnagnstr7hfnwpj32d3mmceum[.]ipfs[.]example[.]io/
* hXXps://bafybeiaf5rcjjdgyffxkcpniqqsk62dx4wyibizfkoohmmkrxupruxjxq4[.]ipfs[.]example[.]io/
* hXXps://bafybeibnufccmhulch2y53ti3ih63ljv564ruxf5qshn33xisna7hhmubu[.]ipfs[.]example[.]io/
* hXXps://example[.]io/ipfs/QmRQmgRkVHfgoK8GtgnfqG61LHTv46HWJWbLzgrHbK9sJC
* hXXps://example[.]io/ipfs/QmVL6csu6yBc8VSjfxmivJHhtuqQFhdJsiPRyWKQRGo19N
* hXXps://bafkreicbq5kqxhlduni4kasyhs2hpug6ft32uzbx7vtc45er7nwwuntytu[.]ipfs[.]example[.]io/
* hXXps://bafkreidgydpt5w246o6zxwfggsges57zedk6z3nzlupbhz2zf3weppbeo4[.]ipfs[.]example[.]io/
* hXXps://bafybeifltunwld5wwwxqpxbpyepsdh5knfknrbyfkcqnjex7a4haslezly[.]ipfs[.]example[.]io/33runk.html
* hXXps://bafybeidh3wdcpsqif5e33rgmpsv55ddzsbmoretfb6beocz24c75r6czyu[.]ipfs[.]example[.]io/

Please find more information on the link below:
-----------------

8336276,https://example.io/ipfs/QmSgSrFykcvstRFiKPtRvRRYUnptP8MFBNb3a3Uc4rTkwU,http://www.phishtank.com/phish_detail.php?phish_id=8336276,2023-10-18T20:56:57+00:00,yes,2023-10-18T21:33:32+00:00,yes,Other
8336277,https://bafybeicaqd2fdx5uvwqcnuwznp4wvwupdfnagnstr7hfnwpj32d3mmceum.ipfs.example.io/,http://www.phishtank.com/phish_detail.php?phish_id=8336277,2023-10-18T20:56:57+00:00,yes,2023-10-18T21:33:32+00:00,yes,Other
`

	var cids sort.StringSlice = FindCIDs(texts)
	cids.Sort()

	expected := `
QmNjmhpMNg6NdPdEj3BRdPBgEH6pkcNxBeW3vQE2dHgCai
QmRQmgRkVHfgoK8GtgnfqG61LHTv46HWJWbLzgrHbK9sJC
QmSgSrFykcvstRFiKPtRvRRYUnptP8MFBNb3a3Uc4rTkwU
QmSkSu4R4EAvR6HRk7BMjUnH4mGfk89uvyAaRSYiDfk3qW
QmSnoKWrdapdiuWxxq9iD3jG1XiYvHSNmqr7oH48ieYNm9
QmTo78zLFeRwAV9TMHPR8AvqWSjMJrkoBd1w2oTWoyXUTA
QmVFksWknMR58UiAtiB1tGZbVYN2umfKgEPADt2zuqdbDQ
QmVL6csu6yBc8VSjfxmivJHhtuqQFhdJsiPRyWKQRGo19N
QmYmCWjc82UNBn792EnyNsGBuTsLJRtHZfqjGxg2J4w5Nr
QmZtZNXzYMWDymUYEUdYQnDWpYGczes7w1nx9HPtg8mgih
`
	got := "\n" + strings.Join(cids, "\n") + "\n"

	if expected != got {
		t.Log("expected:")
		t.Log(expected)
		t.Log("got:")
		t.Log(got)
		t.Fatal("failed to parse cids correctly")
	}
}
