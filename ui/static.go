// Code generated by "esc -o ui/static.go -prefix ui/static -pkg ui ui/static"; DO NOT EDIT.

package ui

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
)

type _escLocalFS struct{}

var _escLocal _escLocalFS

type _escStaticFS struct{}

var _escStatic _escStaticFS

type _escDirectory struct {
	fs   http.FileSystem
	name string
}

type _escFile struct {
	compressed string
	size       int64
	modtime    int64
	local      string
	isDir      bool

	once sync.Once
	data []byte
	name string
}

func (_escLocalFS) Open(name string) (http.File, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	return os.Open(f.local)
}

func (_escStaticFS) prepare(name string) (*_escFile, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	var err error
	f.once.Do(func() {
		f.name = path.Base(name)
		if f.size == 0 {
			return
		}
		var gr *gzip.Reader
		b64 := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(f.compressed))
		gr, err = gzip.NewReader(b64)
		if err != nil {
			return
		}
		f.data, err = ioutil.ReadAll(gr)
	})
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (fs _escStaticFS) Open(name string) (http.File, error) {
	f, err := fs.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.File()
}

func (dir _escDirectory) Open(name string) (http.File, error) {
	return dir.fs.Open(dir.name + name)
}

func (f *_escFile) File() (http.File, error) {
	type httpFile struct {
		*bytes.Reader
		*_escFile
	}
	return &httpFile{
		Reader:   bytes.NewReader(f.data),
		_escFile: f,
	}, nil
}

func (f *_escFile) Close() error {
	return nil
}

func (f *_escFile) Readdir(count int) ([]os.FileInfo, error) {
	if !f.isDir {
		return nil, fmt.Errorf(" escFile.Readdir: '%s' is not directory", f.name)
	}

	fis, ok := _escDirs[f.local]
	if !ok {
		return nil, fmt.Errorf(" escFile.Readdir: '%s' is directory, but we have no info about content of this dir, local=%s", f.name, f.local)
	}
	limit := count
	if count <= 0 || limit > len(fis) {
		limit = len(fis)
	}

	if len(fis) == 0 && count > 0 {
		return nil, io.EOF
	}

	return fis[0:limit], nil
}

func (f *_escFile) Stat() (os.FileInfo, error) {
	return f, nil
}

func (f *_escFile) Name() string {
	return f.name
}

func (f *_escFile) Size() int64 {
	return f.size
}

func (f *_escFile) Mode() os.FileMode {
	return 0
}

func (f *_escFile) ModTime() time.Time {
	return time.Unix(f.modtime, 0)
}

func (f *_escFile) IsDir() bool {
	return f.isDir
}

func (f *_escFile) Sys() interface{} {
	return f
}

// FS returns a http.Filesystem for the embedded assets. If useLocal is true,
// the filesystem's contents are instead used.
func FS(useLocal bool) http.FileSystem {
	if useLocal {
		return _escLocal
	}
	return _escStatic
}

// Dir returns a http.Filesystem for the embedded assets on a given prefix dir.
// If useLocal is true, the filesystem's contents are instead used.
func Dir(useLocal bool, name string) http.FileSystem {
	if useLocal {
		return _escDirectory{fs: _escLocal, name: name}
	}
	return _escDirectory{fs: _escStatic, name: name}
}

// FSByte returns the named file from the embedded assets. If useLocal is
// true, the filesystem's contents are instead used.
func FSByte(useLocal bool, name string) ([]byte, error) {
	if useLocal {
		f, err := _escLocal.Open(name)
		if err != nil {
			return nil, err
		}
		b, err := ioutil.ReadAll(f)
		_ = f.Close()
		return b, err
	}
	f, err := _escStatic.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.data, nil
}

// FSMustByte is the same as FSByte, but panics if name is not present.
func FSMustByte(useLocal bool, name string) []byte {
	b, err := FSByte(useLocal, name)
	if err != nil {
		panic(err)
	}
	return b
}

// FSString is the string version of FSByte.
func FSString(useLocal bool, name string) (string, error) {
	b, err := FSByte(useLocal, name)
	return string(b), err
}

// FSMustString is the string version of FSMustByte.
func FSMustString(useLocal bool, name string) string {
	return string(FSMustByte(useLocal, name))
}

var _escData = map[string]*_escFile{

	"/deck_list.js": {
		name:    "deck_list.js",
		local:   "ui/static/deck_list.js",
		size:    1463,
		modtime: 1563430695,
		compressed: `
H4sIAAAAAAAC/3RUUYvjNhB+z68Y1IWzuViRvT6WvbXzUK4lhcvTlYNSyq5WnsRqFClIirPdkP9eJNvJ
ZtvDkEgz33zzzWgkfNkZ66HBFd8rD0Jx5+ALis1X6TwcJwDCaOftXnhjkzRaAHwrHX1EBTU0Ruy3qD0V
FrnHXxSGXUL2iqQPV2AayQMv5U2TkAbFJlPS+R54mkwA1ugBB4oxmUW/t/pMc8Y6DLrFxiXx91pbNAV5
4f+tDou6Qfubx61LLnmvze+qpFJrtIvfl1+hfksfQQDUGeuThE/hOYV6DpxqvkWYw3NcpCNsy3fJsAZI
jhCcUxDcNu7RIm/+mYLGF/9osZN4eOQeToHvHALwVCk5n1QcWourmvx0cwwcJwJGCyXFpiZ8t6POc+u/
oXPS6OTDgPkwvTm+SXVKH8a+rrhy+EDmA7Ca8fmkamQXElfCNDi/OY6tC8K+ee5d8k72Ab5wj8m1/jQ9
VbPIMKlmkbHi/YTVxAWWTEm9IWM1vekHNbVmSHwu6L8VBMGuW8PLVmlXk9b73efZ7HA40MMtNXY9Kxhj
M9etCRxk49uaFCWBFuW69f06CP/ZvNSEAYOihKIkff+rHfctNDVZ5vdwu/gkspzmwLIC6H1WQNHlpWCQ
05zeQxG+Ni9FhECRBVtWfP8kWIjKQkT4Xpf3kN8t7rrsri26u9dtCazNiu9hl7Nx22VlW3TlK5m9V8KA
tUXZFeWCvRJYSaVqoo3GHlmFQkPjw3HOlJw/DYN0nse/jdQJIVeX4HLAbyZyvA9yBb0Z5sDSsfvR8jCJ
iPhWQCNXK6hjKGSX4Rieg8ASERXkjF14CO+4VPxZIWhzIAOjCpd8IOvjw+7PoQbCCHyEhq7R9ymmV44k
epZG+zZJ4SPkF/857te9Un8gt+9jo3Nh9tb9r2cp9d6jS/pu/hUvtzDbndGofXgFzhvqlBSYZEWaDjUN
BTeDh03hNh1Og5Kgk0Cfp/efnZ/HszpN/g0AAP//jQlJmbcFAAA=
`,
	},

	"/index.html": {
		name:    "index.html",
		local:   "ui/static/index.html",
		size:    455,
		modtime: 1563585442,
		compressed: `
H4sIAAAAAAAC/2xRsU4rMRDs7yv2bf3CKR1C9kmIUENBQ2nsDV7i8528mwv5e+RzElFQeUazO5rxmn+7
l6e399dniDqmoetMfSG5/GmRMg4dgInkQgUAZiR14KMrQmrxqPvNPUL/W8xuJIsL02meiiL4KStltXji
oNEGWtjTZiX/gTMru7QR7xLZ7WrVvJQ10bAjf4DHeTZ9401LnA9QKFkUPSeSSKQIsdDe4ug433kRBD3P
ZFHpW/uVrylN38pU+DGF88Ux8AIcLAbyB0HwyYlYrNEdZyrYxupXbNdM8mD6uL0s94GXa+yrkZAITxmH
v1V1KjetWrQsFYovPCtI8ZcyX7cu4xSOiepeGxo607ez/QQAAP//fSD4gMcBAAA=
`,
	},

	"/main.css": {
		name:    "main.css",
		local:   "ui/static/main.css",
		size:    2243,
		modtime: 1563585362,
		compressed: `
H4sIAAAAAAAC/5RW3W6zOBC991OM1JuvK4yAQpqSm32TlYOH4K2xWds0ya7y7isbaPhJ+6lSpQbPmR/O
nBnTuFbCfwSg1srRmrVCXkuwTFlq0Yj6MJms+BdLyNLuciA3Qo6aX4Nfyy70LLhrSnhNEm/1Z+YkVAkJ
sN5pf9IxzoU6lZDGWYEtJPF+/4btgRCAI6veT0b3ipfwhBXWdX74viApFNIGxalxPmJe+MNKS21KeMqr
YrdjocgmjUiTRaR5iUiTR6QpxpKn8hJIja/iJ9nSYox9p23gJo33Seaj3QiJnXAS14gsLna7EdFkW//d
p3/zsrXmWT5Z8601u0cuttbUsx6stmVSRiR2eHF/hYc1+rM3N0K6iNTatDPaqNNdCcm9zfSondPt1Nrg
JsXcYwIk8e4T0cuIaPmjuDMVzWESa7dI3kuQopTMOqpr6q4dhlSrs4cFBv9Kc/ySk7lc6Sg5czqyX0kE
41+c7J8X1RbdBXbDZBy14WioYVz0toS39bxM46XYRyiBC9tJdi2hlhig/j/lwmDlhFalV33fquDTMqEe
O/3dWyfqK620cqhcCRUqh+Y38VgIth0rVn4IKxzyuV2oBo1wg/h9HiYUmoBYzf5hNfLnRjj8mpujvlDb
MK7PI9HpPo8gfcsjyNLC050Vz5B0F8gHmiMCAF9A0wHpG/I6Mj0r9o9BIFUjJP9aHrFFa4VWP+sPAJPi
pKhw2Np7A+bxPgdtFXSOeRKq611ATRvppRiIGmmmZlxUW6aVVjOi78+L3ddqpW3HKlwm7r4hxKDt5VCU
3yo0vOnyHTlW71QKO6D8D2rdVeK9iMU0J2uvcaE8YOaOYWUtzOMBn3OydOLiYw4c0g+31hIYW8ecpVKo
92UpQoULYhLAwzavl9W2DvtxGnaOkHJ5iY2JLUqsHDX6/Fh43+e9r+DCTMk3ccdrZXY9BrT/sVyMixub
7zniWwA8/dOjddNsBFmdR5UqbVomA+rPFrlgYCuDqIApDr9mHxGF/4h4DgFWqtlqBOBGbuT/AAAA//91
AU91wwgAAA==
`,
	},

	"/main.js": {
		name:    "main.js",
		local:   "ui/static/main.js",
		size:    2767,
		modtime: 1563432096,
		compressed: `
H4sIAAAAAAAC/6xWTW/jNhC9+1dMtT1IgMv0HEMF0i+gRbHbbgr0GHOlScSNLLIcer2Cof9e8EsSLTnJ
oTdhPPP45s3j0OKgpDbwM1bPfwgy8KjlATJ2U2P1/NAKMuwzZbtNSPuIXwSe7pFIyG7M1S76QD6cFNwb
bigBJhuZIW+qlhPBnVJw3gBUsiOjj5WROi9cBMA0glgdGZbQ4WkknBe7KefBJhGUUMvqeMDOsCc0v7Ro
P3/sf6vzzCVkyxrGlcKu/qkRbZ0n5zH09cVuMxXR2JZnM7aZ0nFpL9FxCdmyZklnPHGVj07m4jkls0p4
JdmMjp8Owtx1dEINJfDw8UPkM/s59z9ex3pC8x6/mr+OSMYz4dR3FeSFBfTTDDMGGtnyExcmHPeIpmos
RmS8znmsDV8+d0iEHHOuy+9TLgYQzliMIKUwjQFgsKPQ2NU4mfYkulqemOyUVHZ4COWFDH6ujTxZL1Ps
dwhz9SI1nBooI1grK251ZTa8mx/TCDJS90yjanmF1o+Yn4ctZBb8TqlsC9lNbFQ8Qv6NxSheIWO7Mkcd
5d2M1bbYelIb+keYJs/eBSuvILrL4SsCvZj+nWWVFbFzwJbwot4eET18AZEWxzF4w806mS+RcNnJ9C2y
WpBqeW+vy7Ft1yxwmZd1ssNseVuv5S23F4sbauH5me7zPub929otVFzX9BF53cfW7ECmKJRlCd8X4+TW
bKKO1ASPgEWFYQt7vzBu4duzDQ37Lezfxe+VjfkGca7IuJD7jSKmF9DxLh3967ttdccsRF2oPrrW/Zre
6Bc1dFWphN7pMyH/RyXfakiv+Mr7xeITtfDjrP8LdR40/msXfK64abYglV1JFuM8RKn88tI4IQfxHHRa
ONtJGonJ52jdUGmDn8k/YS6Tt6hNvv+VixZrMBIcqBV9KjDuBRl9GwCj8VwzybULvEOelyG2Of/TMKtc
GmS1eu8UvkkukkdZ2vA1HG1uJmfFPXlA08j6FrI/P9z/nfllOD9l5f1+sVuNJNsvmF07YRuin2Td38Lv
9x/eMzJadE/isc/P8f/DUCyYzJ72Fwl0+NVEtYfNJjiHKxX+2NwpZb0wxVl8eneb/wIAAP//oFoS6c8K
AAA=
`,
	},

	"/review_session.js": {
		name:    "review_session.js",
		local:   "ui/static/review_session.js",
		size:    2440,
		modtime: 1563584857,
		compressed: `
H4sIAAAAAAAC/5xVzY7bNhC+6ykGbFHYwFZuEORSSwIKpEUPadF2e89ypbFNLE0q5MiOsdArtIeitz5d
nqTgj2RR6+waOewP6W++GX7zzbjWyhIQ7lvJCaGEu6zYIW/QVBlAsXtdvcX64XsobMsViKZkDdYPrCpW
7qIqVrvXAfim+s3orUFrp+A23k0D3lRZsRpyZMWeCwW15NaWzKK1QiuotSIuFBoWyF95sg8dWhJasQFO
giQ66t2rKnPAjTZ7DxWq7egnbfaeAKDwF+ePGPCO9EbXnfX/1XrfSiQsmd5s2CoJolOLJbPd/V4QgwOX
HZbs099/MfC4YuWyhvytz2DQdpLGKuMxco7ScGWPaG6JE7LqG3Vv2yjRDFdrY7CmHzz8CbJYtU5OJ2KV
3a2zDD+22hA0uOGdpFAD/IEHgcfbqO5jBuD7brqatFlE1Zf+AwDaCZu/RwklNLru9qgorw1ywh8lutOC
NeLAlusEnQul0Pz85y/voBz9NIN86NCcblGiT8u+OndpmWsVFIYSEMoq1gKAeWvwgIrehhctYl4AsYGF
pxY2iCPUdjnGDWkNWi0PGBDn4B5QWnyC3iL9ih/p9+i0Cd7/7d2xd63eIgFGOYakBqkzanzuiLXo+lE/
LNyvucpPRPHztUzkdFcJW+zY5c4NQ1QOsGkXDKpm0OFM55WPEtVcyns+L3QKgRIGUEIzF+8y0wx1iWys
cho46bJzmOkSd3Vtwwn9NE0el15HtrDyHmFYJzdAmri8AYkbuoEwl+9QQe/yTCVdZ9lgPIeFsizhu7Pl
jkI1+pjvhCVtTrl70tlAwRzh1GfPW2DcmqkN7r5+9JXCt77UfhXP/d0LczYuzpRvuH4herqoUgLmlxF7
IT5dYF/CEFb2Mrd0kpgfRUO7IMfYq77e3V1J4he4S31tVv8tMTEVtydVw3yxJObiw5xcV836SahX+7n4
pCnT+ETt5xhmbVlnnxu1DZc2KTF8oUEJ/MgFXVgQi1BdrMuNSwjJhX0f056nZvKS1Bqf/vuHrS+Aggtq
LbV7H9saRDUCk1d93mqz9f9MDf9eU4PB5ooKoggRMKwC99Nn/wcAAP//fC4cfogJAAA=
`,
	},

	"/stats_list.js": {
		name:    "stats_list.js",
		local:   "ui/static/stats_list.js",
		size:    2212,
		modtime: 1563585394,
		compressed: `
H4sIAAAAAAAC/5RVTW/bRhC981cMti5AtgkZODrZJIsgLpAAziXuzRCqBXckbkPtqsuhVEPRfy/2gxTJ
OFZyEcTZt2/emxkOK61aAsLtruGEUMAqymvkAk0ZAeT12/IOqy83kLc7rkCKgrXEqbVBVuaZjZZ5Vr91
aCH3UDW8bQPqdYsNVvTa6AOzAEu4KN9zI9obyLN6EYIedma/ly05dhd31JmQ+zLKs15blG+5VH22Sivi
UqFhver3nTGoCB4sH6y1mTuwKmYO7NWuCZoa6f9YeWS02pT3vCX4jHuJBxTwjqyFcDQg+xQmwN7RkMPT
Zj3vMwk+KkKz581LxDJgfoL2Tq7XsuoaenqJWAyon6D+IFvSRlYva64H1PPUeWarnme2o2W0uo0i/G+n
DYHANe8a8l32rbSjAccIwA2u6SrSJk5cBIBq2aZ/YwMFCF11W1SUVgY54Z8N2qeYCblnye0EnUql0Hz4
69M9FMObYCEnOxEbJMBwu89jkDqjBoIB2yKBG67Y/U5luRAUHjBWYFAJNB8Jt22cDFzT8JhJYPUFCjhI
JfQh9cV9Si0turOZu387NE8P7k3SJma/nN/fZOq8576NHIHfDL3qkQfPL9fgXaYNqg3VUBQFvElCbQKH
v7XTu87W1IsItfm2Ag/nyj2+WSYTiu/5cJsiSbWqaq42dn/FRyBubNdOCRRlKNzYTt+DdC2ViOMjVNyI
gHZ/rRPPke5502HQ+j214fh0bt2zhmcz+qKfcV9cPAhIt3w3U7zK9Y6kVuXV0UZPeRaeV0l/6R8tVczY
N7PlLXi2V6HTp16nr1e/bqCAT5zq1OhOifh6Ab95/CPrlxZbhjpMLz6QkWoDRdASD4RlAdcL+ANWV0fH
vG60NufjDK4XyUnACm6AsQR+DwSrq+OA+RWuF6d69UOD4tb9tLDW9oV3ZbTGp3cVHuCOE8ZBViiG/UR8
Pt9ZutMkJX2vK96gL0acXMg67PhpzmlJL3CM1vkz4/TIzh8FtrxANVrfU6qZ+fPHgC3h61d4DPZnszu0
z89v/zib1lfQz+sp+j8AAP//7xs2maQIAAA=
`,
	},

	"/": {
		name:  "/",
		local: `ui/static`,
		isDir: true,
	},
}

var _escDirs = map[string][]os.FileInfo{

	"ui/static": {
		_escData["/deck_list.js"],
		_escData["/index.html"],
		_escData["/main.css"],
		_escData["/main.js"],
		_escData["/review_session.js"],
		_escData["/stats_list.js"],
	},
}
