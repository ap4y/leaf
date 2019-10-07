// Code generated by "esc -o ui/static.go -prefix ui/static -pkg ui -ignore tests|node_modules|package.*json|babel.* ui/static"; DO NOT EDIT.

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
		size:    3462,
		modtime: 1570400577,
		compressed: `
H4sIAAAAAAAC/5xXf4/jJhD9n0+BblXptjKWnbVzWUeV+j2qqiJmHNPF4ALeJK32u1eAnfhXstqTTpfY
vHm8DG9m2No2Av+HMK6UtKSiDReXAhsqDTGgebUflgz/Fwq8SdvzHn0gdFDs4uMaeiYnzmxd4B9J4lbd
O33kssAJpp1V7k1LGePyWOA03uTQ4CTe7V6h2SOE8YGWb0etOskK/AQlVFW2fyxIcAmkBn6srWPMcvey
VELpAj9lZb7dUi+yTiNUbyJUv0SoziJU573kQV6CU+1UfGW3NO+5b2kLuUnjXbJxbB8IxZZbAR5h4WwJ
FfwoC1yCtKADwuhyzpAnvwTuzZJ7e+WuX5ar2SYbVrPl6ma7HVbz5WrqTsSvmoYKEaHYSf7LP8zR13P7
QKiNUKV0M0opsaotcHKzADkoa1UzHLsPE3wcMQCSeHtFdCJCSnyJd+SwMUxAZSebdwILXghqLFEVsZcW
/Fazd6sCfXypGNzNydjKpLejPh7o9yTC/b842T1P1ObtGW9D1RyUZqCJpox3psCv81oaSk/Sdy+BcdMK
eilwJcBD3SdhXENpuXJuU6JrpI9pKJfrQX93xvLqQkolLUh78+hDPurJliVHi3duuAU2XueyBs1t6Bud
tUpGiMu2s3+4ZP/2zXSHhttvf/oYcoLDG7eEti1QTWUJBZZKghM0FOFLHnJzO3OcxLl+3E9Cem9ks3S/
9NmNXSIol6C9nFnj2s/4TzW3cP/wDupMTE2ZOvVOSHdZhNPXLMKbNHd+2OTPOGnPOAs+iBDG+A40DUjn
mB9Lsb8GB5c1F+y+f2MDxnAlv2YgjH3/ItxCY2ZdrOe7doIZ6RgTa2q5PJIr9if2876ZMZw0bQvs/n9o
6Fl4HGxHgiM929xf/VwLj7PwJ/99NWw0X2I/7EJPGvtmasLb82QMNUoq09ISZkM4DT68O1j6dLcPbKDB
dMI+nE8MyjciuAko94UYexGjcpw02WQe1ff5FT/cMLSouF7vu7qft8NAvQUx/j4Ghu3DRWMKjI2l1hDB
5dtUCpd+pg+2XzXbfIYsdZj3Y/AgF2J67+g3NiCgtESr03q5Pd73NhlDZ1vl7af9xHF6sNx4Xk2aItsx
gFcPePqnA2OHjuBNdurtLJVuqPj0stKLOmra1tNyCNfBOSSmZ278pcqjjdXqDaaiVvDOpYuryybUW38A
nxOsuC24X5a1G1TGUm1XCLg7mncqcEttPVV9yHZA2aOYkuuyvwr2Qj+Puf3aTyIqWlqlV3TRlwPsyvsR
S1WfRcw13cWHj6ukAA894w52JGb4AQ/xq2bIejOMLXxQgu3np9xwxsQKuXfI4ho4CQUZDuH3Bhin2JQa
QGIqGf4++jsod8Z/9tGzLrrsmRh/oA/0fwAAAP//jgaIKYYNAAA=
`,
	},

	"/main.js": {
		name:    "main.js",
		local:   "ui/static/main.js",
		size:    2810,
		modtime: 1567739819,
		compressed: `
H4sIAAAAAAAC/6xWTY/jNgy9+1ew3h5sIPX0PIELTL+AFkW32CnQY6LazFo7jqWKyqRG4P9e6Mu2Yicz
h94Chnx6fHyizI9SKA0/YvXyGycNByWOkBYPNVYvu5aTLr5Quk182id85Xh+RiIuujFX2eiOXDgqeNZM
UwRMJjJDTqqWEcGTlHBJACrRkVanSguV5TYCoBtORR0YltDheSSc5dspZ2eSCEqoRXU6YqeLz6h/atH8
/L7/pc5Sm5AuawomJXb1Dw1v6yw6r0BXn2+TqYjGthybsc2Yjk27R8cmpMuaJZ3xxFU+KpqL4xTNKuIV
ZRcKSbSv+NTRGRWUkOVQfufJRP/dwWD1K+sqnAgw6rsKqBIKDZqbo58u0JR2Zlz7s2KMzNb6E1fPnFD8
L5c7RGKOObdH4FKuhhD6uh5DTGEaBcBgxqGwq3Ey7pl3tTgXopNCmgHiKO9l3hc14mz8TEHhwc/WydUw
aqAMYK2omDaHm/B2fkzDSQvVFwplyyo0nsTsMmwgNeBPUqYbSB9Co/wA2VcGI3+DjOlKn1SQNxmrTbHx
pdL0F9dNln7wdl5BtBfEVXh6If0bwyrNQ+eALeFVvTki2OIKIi4OY/DumzqZLxJ/4Un3LRY1J9my3lyZ
U9uuWeA6L+1Eh+nyxt7KW26wImypufsPqKtmpvu8j3n/pnYDFVM1fUJW96E1M5ApCmVZwrf5OLk1m8gT
Nd4jYFBh2MDeLY1H+PpiQsN+A/sP4ffK1nyHODdkXMj9ThHjC2h5l5b+zd20vm0Woi5UH11r/41v9F0N
bVUsoXP6TMj/Ucn3GtIpvvKG3VBx+j88Ywu/zvS5Um+n8J8Tks4k080GhDQry2BchiClW24KJ2QvroWO
C2c7SyEV4iVY21ea4Bdyz5zNZC0qne1/ZrzFGrQAC2qGMhVo/Fdn+ehrDxiMaZuJrqXn7fPCA+nanH9Y
zCqXBlqt3luFH6KL5lCWNn0LR+mHyXlhjx5RN6J+hPSPj89/pm5Zzk65eubvNupz562uvtx3QXxFeovh
xkf/FnX/CL8+f/y9IK1495kf+uzivyuGPGpkSBJvICal/wZ6ktJYYooX4YXeJv8FAAD//1mSjsL6CgAA
`,
	},

	"/rater.js": {
		name:    "rater.js",
		local:   "ui/static/rater.js",
		size:    3593,
		modtime: 1568852625,
		compressed: `
H4sIAAAAAAAC/8RWzY7bNhC+6ymmbFFIRVfen5y6loAFsmgCpAG67nGBLE2NbXZpUiEpO8bCr9Aeit76
dHmSQhQli16vLSQFerJJzffNcIbDb5iSxgKtrLqjFgvI4CEaz5ReAi8ywmVZ2bN6SYAJakywlUcAY7fe
GRPHNVOsMu4fU8tSoMWMqNmMjHoQuykxI6aaLrnt2Jvl2bSyVkkCKyoqzMjnP/8gMMqj8aj2m0fRuHQe
NZpK7MB+6XyYkkpnQ6VZoz4zllok+fdyasrxqP4amjGlNTJ71pjvGY5HZR49XEcRfiqVto0/uPFJ0/AU
AbhE6opZpePE7QDYBTepYUojZHB+vdv7gAIyKBSrlihtyjRSi7cC61VMCr4iSWidcilRv/ntl3eQ7aq1
Z/OxQr2ZoEAXBPm2V6okVbLJLWSAkOU+QgBMS40rlPY1zmglbOw9t7xKThwu3h3GW2zrn20UAczRAvro
27NrtJWWXXCdrUELHSejQkwpewwS1rmEDFqDHXyh1r9WaCxXMn6CplwfBMq5XcA25DmQkuA2JEFWiSs5
OZHTvXvyAoXjaFqruezZiSq19XaL1NiNwHTNC7uoG/K7p0v4ITzqli0e+gjXKHUEpL/r2rApaJe8O9ck
ccPW5qsJtTKo3w4Mt/F43UM3jJM6s8fwYQX6BD61N87gGMV+DYJ8l9Ra1LJuE/c5NaXgNh7dm1GS/q64
jMn9venyPYO4O3a6pJYtYolruMP57acy9mRJknT90jtmWPvP//xFrg8YNcVkSqj6UGSuEWVnGJz55eu4
95Rc+AYEFAaHRPb3kMg0FgPiaggOhuVfuK27btsoagpiUMx62jL85a6BR97jnkxpXHFchzqlqeVyflCo
aLGikiHZkyAvNpOFWkNzcAKUMTTmETcZMSVlWCtYTVbwVePZeXEOAMaNbLVE5wH8guQ3c8rleNRYHYRc
BJBLkr+hujiKuAwQVyT/WanjiKsA8Yrkt9Rs+ojxqOCrndg+072JL+gJ3fs6jetujW/ujocWxW0tV++4
sShRx+QRN4VaS/JjKGx1a3/TcZdUo7TvVYFJK02yEsKzA5g1t2wBcd0PtU3LAsCoQSATV/2fut09gesU
rVPPFviaz7m9GIA8Pwy9HAC9OAy9GgC9PAx9NQB61UG3UWDXJNZ9SE7Iab9z/4MZpT+ZvOT1RoiY+AEz
SWdK31K2iL1M5+3lcXrw/K4xwdnj/k07FuCzEN9Xyynq+j4afCttjKmleo5ewpNdThP3J/nfh6yP/s/J
2ar/Xn/JaOXf0sTr0oobPuWC203NsOBF0enmy8OFf9kPc7iVQDJ0Hho8kAQnb0eaUwL6FcnoH+SLs9HL
aK3W/wYAAP//fRdB/QkOAAA=
`,
	},

	"/review_session.js": {
		name:    "review_session.js",
		local:   "ui/static/review_session.js",
		size:    2027,
		modtime: 1567985489,
		compressed: `
H4sIAAAAAAAC/5xVwW7jNhC96yumag8y4MoNglxiS0XQFuhhF9iN9x4z0tjihiYVcmSvYejfF5RISXSc
HPZim8SbNzOPb8Z8XytNcIaHhtQjI9RzWKPYdj+hha1We4jThbbn9LuJl1FUKGkICPe1YISQwSZaVchK
1HkEsKpu83+xeLmHlamZBF5mcYnFS5yvFvYiXy2q2x54l3/RaqfRmCm4dnfTgLs8Wi18jmi1Z7zHGjSG
KxlDIZgxwxkKJYlxiTruM9106NcGDU3hxEmgzVPd2ASWNo82yyjCH50sJW5ZI6iHwyMeOB7XLsU5AuiU
0E1BSiez7gaAKm7SJxSQQamKZo+S0kIjI/xPoD0lcckP8WwZoFMuJer/v33+BNkg7TIaMWZ4lAwkHsdH
SqZMAypVct087zlBBkyao43LO5zPWTFZCuwpwti5i5hN8zPvD5d/8EuQf0D9Qn42OnDID9DaGnZIgE49
r7NGarQcBBywBu27FS+J/bh8lPS1QX1ao8DuzeLfO2fOAvXtVcDmTJW475DTOy7zsOm7apSlV2ig02iU
OOBD12NSMCGe2WWlAQYy8KiAiJUHJgt0hnyHKQRdoxqqnMZx0+fmcmcNqRuc9tXUJSNcEyOcNBdeO7Z+
WZzBz94cSBETcxC4pTloRlzunuhUI7Q201RU5z++hcSiIcsy+MsTAxy5LNUxrbghpU+pbcqb0dujP7XR
xyYYNk5ohM0f565W+LMrtl24c7tZfsw37JmQz1+7tnpltJupqRC2z9iOYwx/X07//cWgLSdcoxffLc0v
TKeTldZdpUXFRalRpgLljqquiptRbQ9jdY2y/MeCk/5fwU2mo2wBhcE3YRprwQq8EjcfMIIZ6gCeqh/0
DmwqdfzqBEwCl4wGZOYkCwg2i+4XiimURt+LbfrS5GOfV+y/ZcKgN5Y3tNtq1rTsyDhdm9yJHYcmHtE0
gpLJinujWVfgqUa17QuH36whZLN/Rh3PAmtfH/Ok73dUsY3a6GcAAAD//50A0bzrBwAA
`,
	},

	"/stats_graph.js": {
		name:    "stats_graph.js",
		local:   "ui/static/stats_graph.js",
		size:    2917,
		modtime: 1567396628,
		compressed: `
H4sIAAAAAAAC/5RW3W7jNhO991MMCH8AlbUl2V9yk1guWmzdvUiKAim6GxRFrUi0zEIWDWpsayHo3Qv+
SKLsjZu9cULO8JwzR+SQrNoLiZCyTXzIEZI8Lkt4xhjLX2S830I9AkhEUaI8JCgk9fQMAG556f/Ncogg
Fclhxwr0E8liZD/nTI1+faZki7i/D4LT6eSf/u8LmQXzMAyD8piRCRD1x3sYoPma/5GX6MdpSkmphEwz
pcSkNqMRQMYQmGHp9EiGB1l0QF1uyRA0CtW/Q/l6CiKT4CqRrEiZ1A7Qnnc4bZG0OdAiObgPTjgRhwJb
Ij9nRYZbE+cboCa6gLlny3gYOWtPPMVtB6094qzAz2ra5dgynm3xMvGTnj/zuWT4I6LkrwdklBw5O/0k
KjKBdQghjGvN2cC4NqDN2htI2sUy4wVEUAOK/T3MwgnkbIP3MLubgFRLzOSrQBQ7NQ2Nq5UXBZOfbWGm
wKlF9RVQP5K9eGfpp7ZWW3SXjmLfDwz7QHmJscQv7Zf4M/zLH34oVqR91HyXKczOs8okzpnKc+oIgOrF
U8vhuQtUrtofS6BY9ilwY6HO3K1eIIKnGLf+Lq6o7/tm3+ziPaU18AKZPMY5NJ5CbIeeBx9gdiHzpZVp
PQsMcMJ4ThXTQKdJtvBKrdYy7TlaxS/nilfXFG9i1TusXjN4Q+3qutrVQO1GybXYVuwKpi1BK3U1kLqP
ncM0EInlxLF2AjytlF69FtQIoiiCEH6A9dO4/kKx9JrJuH6hnTnNGu5h/fhWUCMNC/jtqprWtndr2bxQ
W/ylFDfWKnGkpOK8f2kllvHCHv0l14uEyyRnkFQRaakIJF/VyC2dgIzIHQmW3/Lg43uYBzvoKq9T5yWt
w3uM8wP73pJtDGC9QFYhDPjPy1adI2zIclzn8SvLHT8WgVr9TTf+eJ8s14//FtV70msyEz6KFa9YSueX
qoY3hjmVvz89QgTr0SIDlHFRboTcRUT/m8fI6Lh2urjadn1fbjyipC4y88iISFzxkhj1i5wXDKpZREIC
X2dKs9MFGgLVvJvT7VZVNr9IWy4CBWQxWzNC64LbVz7A7E7bULATfFTKrYv2griBWRiGxgkfxaNQnUTl
PaPkRUY7t4ZU5xLfTauvju8hXQTZcjSwU7+S7MZobdXdLtVbQLUa/x/BC0qI15DAZIxrff6admC2X3OF
oN3FlxT7NxhSh+A4xFc/a/O+akajzaFIkIsCzs5L/wiFrThIdTq6INzA/NbrtrE+TLaj8aN6IOgFAcxv
+8hOpF3kfzZi34/me1ATW0Ywv1U9dlzrW2iTCyFpyo9ek4JqsISoW4wqvKXpxuN6J9Jma4MjfY6a0b8B
AAD//xIbl/xlCwAA
`,
	},

	"/stats_list.js": {
		name:    "stats_list.js",
		local:   "ui/static/stats_list.js",
		size:    2147,
		modtime: 1568850503,
		compressed: `
H4sIAAAAAAAC/5RVTW/bRhC981cM1i5AtjUZODrZJIvALpoAzqXuzTCiBXckbkztssuh3EDQfy/2gxQp
K45zEcjZt2/mzYwe5abVhuCeOHV/Gd7WsDJ6AyzNOhv6srax9GvHrqOo0qojINy0DSeEApZRXiMXaMoI
IK/fl7dYPV1B3rVcgRQFcxwXAqsnVuaZDZd5Vr93cCG3UDW86wZYhw1WdGH0M7MAy7gob7gR3RXkWb0I
QQ+b0DeyI0fvDhx3JuS2jPJsqC7KN1yqIV2lFXGp0LCh7pveGFShDbDS5oWGihtxpMHe7ZtQVSP9gy2Q
jFbr8o53BH/jVuIzCvhAVkQ4GpFDDhNgF5zGJJ43G4hPZPikCM2WN68xy4D5CdpbuVrJqm/o22vEYkSd
ps4z25s8s40vo+V1FOF/btUErnjfkB+G7/id7Ah2EYDbMNNXpE2cuAgA1bJLv2ADBQhd9RtUlFYGOeGf
Ddq3mAm5Zcn1DJ1KpdB8/OfzHRTjyk4hbrGhAIXPk+2Pj2n+7dF8u3ebpU3MrBqWpLxtUYmbWjYintCl
6CtyJHu7IGskCMFRkEHqjRpTjNgObXOqp9j+HKs/LuRs8tdKZmJtaMbpkLH7nbO6EBQeMNVtUAk0nwg3
XXzQMg8HJu8JA9OE1/PJFfjMaYNqTTUURQHvktCD62iStNVtb4fkNYZ6X1Z1f1Dz8O4xmVF8t03OIpJU
q6rmam2tK94BcWPHs0+gKIOaqZ6hMelKKhHHO7AWENDu0UrxHOmWNz2GYr9XbjjeH/p5UvHb5h4ETefu
DkIF6Ya3RyUvc92S1Ko839noPs/C+zIZLn3VUsWMvZi41+DZfg+z3s8XYLAZKOAzpzo1ulcivlzArx7/
wAazYo+hEfOL92SkWkMRaolHwrKAywX8AcvznWNeNVqbw3EGl4tkL2AJV8BYAr8FguX5bsT8ApeLfb18
26o4o5931oZet4WzqYHPL1t/ueWEcSgstMN+HYaPwwdij+40SUnf6Yo36NvxIzc6G919nnPe1B9wTIz8
xEI9sMPngD2+MNB05iEP7KPsSBtZ2Tn7NdpH/wcAAP//jj7YDGMIAAA=
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
		_escData["/rater.js"],
		_escData["/review_session.js"],
		_escData["/stats_graph.js"],
		_escData["/stats_list.js"],
	},
}
