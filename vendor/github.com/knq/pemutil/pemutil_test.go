package pemutil

import (
	"crypto/elliptic"
	"io/ioutil"
	"path"
	"sort"
	"strings"
	"testing"
)

func TestBlockTypeString(t *testing.T) {
	if PrivateKey.String() != "PRIVATE KEY" {
		t.Errorf("PrivateKey should be PRIVATE KEY")
	}
}

func TestPEMErrors(t *testing.T) {
	tests := []string{
		// bad bytes
		"----",

		// bad rsa private key
		`-----BEGIN RSA PRIVATE KEY-----
miieowibaakcaqea4f5wg5l2hkstenem/v41fgnjm6godrj8ym3rfkeu/wt8rdtn
sgfezoqphegq7jl38xufu0y3g6ayw9qt0hj7mcpz9er5qlamxjwzxzhzaahlfa0i
cqabvjomvqtzd6uqv6wpeyztdtwiqi9axwbphsspnpygin20zzunlx2brclcihhc
puiizoqn/mmqtd31jsyjoqov7mhhmtatkjx2xrhhr+1dckjzqbstagnpyvaqpsar
ap+nwripr3nutuxygohbtsmjj2usseqxhi3bodire1autyhceabewn8b462yewka
rdpd9ajqw5sivpfdsz5b6glyq5ldyktzntuy7widaqabaoibaqcwia1k7+2oz2d3
n6agCAbqIE1QXfCmh41ZqJHbOY3oRQG3X1wpcGH4Gk+O+zDVTV2JszdcOt7E5dAy
MaomETAhRxB7hlIOnEN7WKm+dGNrKRvV0wDU5ReFMRHg31/Lnu8c+5BvGjZX+ky9
POIhFFYJqwCRlopGSUIxmVj5rSgtzk3iWOQXr+ah1bjEXvlxDOWkHN6YfpV5ThdE
KdBIPGEVqa63r9n2h+qazKrtiRqJqGnOrHzOECYbRFYhexsNFz7YT02xdfSHn7gM
IvabDDP/Qp0PjE1jdouiMaFHYnLBbgvlnZW9yuVf/rpXTUq/njxIXMmvmEyyvSDn
FcFikB8pAoGBAPF77hK4m3/rdGT7X8a/gwvZ2R121aBcdPwEaUhvj/36dx596zvY
mEOjrWfZhF083/nYWE2kVquj2wjs+otCLfifEEgXcVPTnEOPO9Zg3uNSL0nNQghj
FuD3iGLTUBCtM66oTe0jLSslHe8gLGEQqyMzHOzYxNqibxcOZIe8Qt0NAoGBAO+U
I5+XWjWEgDmvyC3TrOSf/KCGjtu0TSv30ipv27bDLMrpvPmD/5lpptTFwcxvVhCs
2b+chCjlghFSWFbBULBrfci2FtliClOVMYrlNBdUSJhf3aYSG2Doe6Bgt1n2CpNn
/iu37Y3NfemZBJA7hNl4dYe+f+uzM87cdQ214+jrAoGAXA0XxX8ll2+ToOLJsaNT
OvNB9h9Uc5qK5X5w+7G7O998BN2PC/MWp8H+2fVqpXgNENpNXttkRm1hk1dych86
EunfdPuqsX+as44oCyJGFHVBnWpm33eWQw9YqANRI+pCJzP08I5WK3osnPiwshd+
hR54yjgfYhBFNI7B95PmEQkCgYBzFSz7h1+s34Ycr8SvxsOBWxymG5zaCsUbPsL0
4aCgLScCHb9J+E86aVbbVFdglYa5Id7DPTL61ixhl7WZjujspeXZGSbmq0Kcnckb
mDgqkLECiOJW2NHP/j0McAkDLL4tysF8TLDO8gvuvzNC+WQ6drO2ThrypLVZQ+ry
eBIPmwKBgEZxhqa0gVvHQG/7Od69KWj4eJP28kq13RhKay8JOoN0vPmspXJo1HY3
CKuHRG+AP579dncdUnOMvfXOtkdM4vk0+hWASBQzM9xzVcztCa+koAugjVaLS9A+
9uQoqEeVNTckxx0S2bYevRy7hGQmUJTyQm3j1zEUR5jpdbL83Fbq
-----END RSA PRIVATE KEY-----`,

		// bad ec private key
		`-----BEGIN EC PRIVATE KEY-----
mhccaqeeiah5qa3rmqqquu0vbkv/+zouz/y/iy2plpicwusyimswoaogccqgsm49
awehouqdqgaeyd54v/vp+54p9dxaryqx4mpcm+hkriqznasysorqhq/6s6ps8tpm
cT+KvIIC8W/e9k0W7Cm72M1P9jU7SLf/vg==
-----END EC PRIVATE KEY-----`,

		// bad certificate
		`-----BEGIN CERTIFICATE-----
miidxtccaq2gawibagibadanbgkqhkig9w0baqsfadcbgzelmakga1uebhmcvvmx
edaobgnvbagtb0fyaxpvbmexezarbgnvbactclnjb3r0c2rhbguxgjaybgnvbaot
eudvrgfkzhkuy29tlcbjbmmumtewlwydvqqdeyhhbybeywrkesbsb290ienlcnrp
zmljyxrlief1dghvcml0esatiecymb4xdta5mdkwmtawmdawmfoxdtm3mtizmtiz
ntk1ovowgymxczajbgnvbaytalvtmrawdgydvqqiewdbcml6b25hmrmweqydvqqh
ewpty290dhnkywxlmrowgaydvqqkexfhb0rhzgr5lmnvbswgsw5jljexmc8ga1ue
axmor28grgfkzhkgum9vdcbdzxj0awzpy2f0zsbbdxrob3jpdhkglsbhmjccasiw
dqyjkozihvcnaqebbqadggepadccaqocggebal9xygjx+lk09xvjgkp3gely6skd
e6bfiembo4tx5ovjnyfq9oqbtqc023cyxzibsqu+b07u9pppl1kwiuergvzr4oah
/PMWdYA5UXvl+TW2dE6pjYIT5LY/qQOD+qK+ihVqf94Lw7YZFAXK6sOoBJQ7Rnwy
DfMAZiLIjWltNowRGLfTshxgtDj6AozO091GB94KPutdfMh8+7ArU6SSYmlRJQVh
GkSBjCypQ5Yj36w6gZoOKcUcqeldHraenjAKOc7xiID7S13MMuyFYkMlNAJWJwGR
tDtwKj9useiciAF9n9T521NtYJ2/LOdYq7hfRvzOxBsDPAnrSTFcaUaz4EcCAwEA
AaNCMEAwDwYDVR0TAQH/BAUwAwEB/zAOBgNVHQ8BAf8EBAMCAQYwHQYDVR0OBBYE
FDqahQcQZyi27/a9BUFuIMGU2g/eMA0GCSqGSIb3DQEBCwUAA4IBAQCZ21151fmX
WWcDYfF+OwYxdS2hII5PZYe096acvNjpL9DbWu7PdIxztDhC2gV7+AJ1uP2lsdeu
9tfeE8tTEH6KRtGX+rcuKxGrkLAngPnon1rpN5+r5N9ss4UXnT3ZJE95kTXWXwTr
gIOrmgIttRD02JDHBHNA7XIloKmf7J6raBKZV8aPEjoJpL1E/QYVN8Gb5DKj7Tjo
2GTzLH4U/ALqn83/B2gX2yKQOC16jdFU8WnjXzPKej17CuPKf1855eJ1usV2GDPO
LPAvTK33sefOT6jEm0pUBsV/fdUID+Ic/n4XuKxe9tQWskMJDE32p2u0mYRlynqI
4uJEvlz36hz1
-----END CERTIFICATE-----`,

		// bad block type
		`-----BEGIN HEADERS-----
Header: 1

-----END HEADERS-----`,

		// empty
		"",
	}

	for i, test := range tests {
		s := Store{}
		err := s.Decode([]byte(test))
		if err == nil {
			t.Errorf("test %d expected error, got nil", i)
		}
	}
}

func testPEM(i int, name string, exp []BlockType, t *testing.T) {
	filepath := "testdata/" + name

	// build PEM
	s := Store{}
	err := s.LoadFile(filepath)
	if err != nil {
		t.Errorf("test %d (%s) expected no error, got: %v", i, filepath, err)
		return
	}

	// check that store len is same as exp len
	if len(exp) != len(s) {
		t.Errorf("test %d (%s) expected length should be %d, got: %d", i, filepath, len(exp), len(s))
		return
	}

	// make sure that all the types are there
	for _, bt := range exp {
		if _, ok := s[bt]; !ok {
			t.Errorf("test %d (%s) should have %s, but not present", i, filepath, bt)
		}
	}
}

func getExpBlockType(suffix string, priv BlockType, pub BlockType) []BlockType {
	switch suffix {
	case "private":
		return []BlockType{priv}
	case "public":
		return []BlockType{pub}
	}
	return []BlockType{priv, pub}
}

func TestTestdata(t *testing.T) {
	files, err := ioutil.ReadDir("testdata")
	if err != nil {
		t.Fatalf("could not load testdata: %v", err)
	}

	for i, f := range files {
		fn := f.Name()
		if strings.HasSuffix(fn, ".pem") {
			base := strings.TrimSuffix(path.Base(fn), ".pem")

			// get key suffix
			var suffix = ""
			if s := strings.Split(base, "-"); len(s) > 1 {
				suffix = s[1]
			}

			// get expected block types
			var test []BlockType
			switch base[:1] {
			case "b": // base64
				test = getExpBlockType(suffix, PrivateKey, PublicKey)
			case "e": // ec
				test = getExpBlockType(suffix, ECPrivateKey, PublicKey)
			case "r": // rsa pkcs1
				test = getExpBlockType(suffix, RSAPrivateKey, PublicKey)
			case "p": // rsa pkcs8
				test = getExpBlockType(suffix, RSAPrivateKey, PublicKey)
			case "c": // certificate
				test = []BlockType{Certificate}
			}

			testPEM(i, fn, test, t)
		}
	}
}

func TestGenKeys(t *testing.T) {
	symStore, err := GenerateSymmetricKeySet(256)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	rsaStore, err := GenerateRSAKeySet(2048)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	ecStore, err := GenerateECKeySet(elliptic.P521())
	if err != nil {
		t.Fatalf("expcted no error, got: %v", err)
	}

	for i, s := range []Store{symStore, rsaStore, ecStore} {
		// serialize
		buf, err := s.Bytes()
		if err != nil {
			t.Errorf("test %d expected no error, got: %v", i, err)
			continue
		}
		if len(buf) == 0 {
			t.Errorf("test %d buf len should not be 0", i)
			continue
		}

		// unserialize
		s0 := make(Store)
		err = s0.Decode(buf)
		if err != nil {
			t.Errorf("test %d expected no error, got: %v", i, err)
			continue
		}
		if len(s) != len(s0) {
			t.Errorf("test %d s should have same length as s0 after load (%d!=%d)", i, len(s), len(s0))
			continue
		}

		// check that the same keys present
		sKeys := keys(s)
		s0Keys := keys(s0)
		if len(sKeys) != len(s0Keys) {
			t.Errorf("test %d sKeys and s0Keys should have same len", i)
			continue
		}
		for j, k := range sKeys {
			if k != s0Keys[j] {
				t.Errorf("test %d key %d should be %s, got: %s", i, j, k, s0Keys[j])
				continue
			}
		}
	}
}

type BlockTypeKeys []BlockType

func (btk BlockTypeKeys) Len() int {
	return len(btk)
}

func (btk BlockTypeKeys) Swap(i, j int) {
	btk[i], btk[j] = btk[j], btk[i]
}

func (btk BlockTypeKeys) Less(i, j int) bool {
	return strings.Compare(btk[i].String(), btk[j].String()) < 0
}

func keys(s Store) BlockTypeKeys {
	k := make(BlockTypeKeys, len(s))
	i := 0
	for key := range s {
		k[i] = key
		i++
	}
	sort.Sort(k)

	return k
}
