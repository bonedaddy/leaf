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

	"/index.html": {
		name:    "index.html",
		local:   "ui/static/index.html",
		size:    959,
		modtime: 1556597664,
		compressed: `
H4sIAAAAAAAC/3RTTY7UPBDd9ynqq8W3G6IRGzRyIo2EWLFA4gRup9Ix7djGVWnoS8AdOQmynb/RDKu4
Xur5vfqx+q8PRu6RYJTJdaeTyl9w2l9aJI/dCUCNpPt8AFATiQYz6sQkLc4yPHzA7lT/iRVH3UcyV3iO
UTU1rv+c9VdI5FpkuTvikUgQxkRDi5O2/p1hRshGWhT6KU2Oi3hT1fPxHPo7GKeZW7wkweVur29g+xZ7
MlfGNcEEL9p6SktaruOxuOMn1YyPGzq7jf3ZsmwXZODBZaRTzewWscbr21pwb6swE7MN/iBEuqe0hhl4
vwcA2cQTKI7ab9JZJAPHPHVOx/BLCpdEzEdqXLDXdNXsorWJ2dEK5J6vhS724XXHSs+KzveZWHKNK6kM
N6vmTu7pQ0hTIVgfZ/kU0oQvKirwnoCgZwlDMDOXkwlTdCTUYhgGbN6g1hXh+TxZQbhpN1OLf37/QmiO
tWcfR1+xaCbi2e0TXsIXKltntecflL6KFsLuf3/m+MaAtmwTUiIjz4X0j3zVxH0geQDrSvX2Vlc9L3hd
dTbJRgFOZnkf3+qIC9ydVFOf698AAAD//0tyxny/AwAA
`,
	},

	"/main.css": {
		name:    "main.css",
		local:   "ui/static/main.css",
		size:    2854,
		modtime: 1556607882,
		compressed: `
H4sIAAAAAAAC/7RW3W6jOBS+Lk9xpN5MpYAggbSlN7vaJzH4kHhrbMY2U7pV3n1lA8UGOrN7McpF4Pjz
d34/m+SiDHxEdy0Z4jdGzbWExzTthpforpHCxA1pGX8vQROhY42KNfOKZv9gCcfMYTkTGF+RXa6mhFPh
bB2hlIlLCcejfb9Fztc1s+48hvy0ZTg/OltL1IWJuJLGyHbmnXgSwwzHFVdRbLmeil9yXY8rnlO+5Sly
n8fIroRiN8ww3dOK+vi4pc631KfdqLPcp87X1F/3wqHeJlslOd1yP3vU3QHcf8+nBzk/VFzWr997aXAy
dGp+SghHZeYXsUCSSg7zYyOlsUt6HLqve+K59qAcG7MCcrZDFRRqidkip7GcqMZm3VVSUVSxNu8cS9CS
M7pYa8mlKkFdKvItPbhfkhUPC2ASTup+fh11X01Z6L7znFts5uFqSbf1fK3oqsHZ82qPXa9I/XpRshd0
N9D06cHXYtENcPZTVoSyXk/tn8q4Dq9TGAQ/Vvd3FW3Ke5vXfdM0/ymX2TjE+kqofBu9nLsB0jCg3AW0
n/RPxzkoRyD4cML95u0I/5T6J+Wo/HAinSnrBqDSGKQ79dz67XYEEcTY80VePV/ELXkgux3RJ8oG/nk+
eAtWTp/2sXZQcqJNXF8Zp34d9+y2pjvmJSlvcSe91Ob2R4uUESCcAxEUvrVMfN5pj2k3PNiN84W3tG86
Im/RLYoqSd/hIwJYZoL0Rr5EAPMQ5n8V5/Of1rJMaAn3WGPT5DaMSJAfjoMy3XHyXkLDcbAb7H9MmcLa
MClKS9m3wu1pCRP7m/7utWHNe1xLYVCYEmoUBtUv+Ah8rCOGW0TKH0wzO0cfS0ZMXFEx47Yl1gthApVD
+MdVcs6esF3n/XZlBp1x5zix1kWDbnSzp/wA2XN+gGNWHCBNjsUDpN0A+SjmwxewbERZyT+OkxwlGrVm
Uvy/WgMQzi4iZgZbvRTT52ukavdIfcw9E11vHCq8aj+LFqvRnG2LJqTwara8B19crRRSd6TG0bFC3fPR
ocHBxC6LMH6K9WvMmTbhBL8EjbRvFjMf26PvcLu7UnfTXzDEH5Ag11sUTXfTSiL0iSI+O8D99x61cc2D
4NtESNUSbmd1lrOuFaKYFL18pRbppGiAbeQ/a/8tCreQaccqlV5oNCP+Fv0bAAD//zOba+QmCwAA
`,
	},

	"/main.js": {
		name:    "main.js",
		local:   "ui/static/main.js",
		size:    4541,
		modtime: 1556598149,
		compressed: `
H4sIAAAAAAAC/8RY247bNhO+91MMmMUfCXHkzf/frQ/A/kiDtMih7eamKAqblsYrZmlSJal1DENPUfSu
T9cnKXjQcb1eNze5kzgHDr+Z+ThSyqnW8BrTu3dMGziMAKjeixQUigxVFLslgFQKbSDD9E7DHOiOMgMm
ZzpZbtCkuXWgo3jqdDOZllsUJrlF8x1H+/j//fdZRLKwDYkTJgSqt5/ev4O59+osARItlYkiOoZ1DPMF
0ETQLcIC1u4hrtW2tIjCM0B0ACscQ0pVppcKabYfg8AvZqnwnuFuSQ1U1l9jArCacebfZxRyhZs5eXZx
sH4qAlKknKV3c0KLItGGKnODWjMpoudB53k8BYWmVAI2lGuckkWQzCY0+M3Y/WKWygwXFwcPlg/nxlCj
o0GwO3hNDUb9qOO4mk2ch9nEehvNJpwtVuEYDRqfJRMRIQ7/atTksJebkMgQ846JTO4Sp+ATo0mcmBxF
pFBb5A/BOduAXUrkXVwb29fPWgqb8KDVkRj8YqLgyz67NHJUJlq9oYxjBkaC29hn/gouDlavWsWhgKr2
ID3EOomtj+Ojo9keFnDZCZBm+xBbqFy22cDcmcLLFu2wofXiNGbw6vKy9UPoPWWcrjmCkDsSPHI0kAVn
oeJhDr8GJMglgReQ2eL3W4x7gshJ3kth8iiGF/CqlTd2b0rOf0GqhrZO+FaWSh+VvGeiNKgjXxa/uR5J
5baQAoXLQvOSaM5SjF7+N64zGA6cBcnlGP4Xh7JKiI2TgN/HyxvhVV101WjkueRnnzDfLi5LLgOqTI1s
+cS1A9PXQu9QMXELczCqxGkrtLVhQS4576zq4LcVdOp9wFmP8hATRWneSLUlcSKFLtdbZmAO2C17TAqF
9yjMa9zQkpu6VHyxDMOPG7sQ51KhlvwevUZrXAFyjQ+0bSn9VKI2zHdVUPbdEFLURSW0b860kWpvGcqg
k53BwAP2bdweRbnL9D0ebMzijt2yLDJq0HZr6C3fxL3l3o1ygN/DscdgpKF8DBw3ZgzUAfcOBVR1kCGq
AIdNg1WF+XxuW7+GdADNmqZ3LaK+zAPLjE5jVSh5q1DrAV6ri4MLFF66SKtJeK9WT2BfH3Tgr16enlGz
JE602XNMdiwzuQ+mAapK86dC8LouD4MoyH/EWhfkCftUKoWp8UX9VR7qY9xTXqK1OtdgI9NSRw9uuH7j
dMnleA27e+eDu6GmZxLRI0UdAhj0ea+4Pdx2wjkLjukDU7fnKfteQrv2vUyd8jBIaZdr+qj4KaezhUJd
cjPkCEemAQ0fXeeO9SYJ08uwbdu2nZP0y+rvv/4g0yNKvhFSyaU9H7lViKJR7J3q8TIdsPGJGP48JwaF
2RkRBBCCQs1F/arqUa1j2WPjWz3++Slu5YwmFwerX63GnUtmiyaX2RWQHz/efCL17dKMj8fGvnMHv38/
+rkw3QgWWvThANgbAbuzrO/cU6OsJYRvMMm6wbKm8pMTbZ3hI43ydI5JoBvyeHbHzfpaZvsr+OHm44dE
G9vFbLOPDjUrNen/ZnXg5y4fzqkaaEbL66JwcQ3mvFAAUhSy0IExo7h7CH8h5HLX+07tDVYPpJ7kcqrz
dt7iMqU2wYldno6OjBsKC07TcFkcqjEQ6/S6KMgYyIR0uNC6iMPe3W6367WbiDyzdiRuy6cT51Njbv1V
51kqY7rgdN8bqR81DZ15xJgIKZBMBz8F3M+Duf+uCq/N34DwntRZa4/yGMkNMC1KnQdAnTeoxrDynxm2
bGq6Wz2rn8/5D6GfONtXQNPi6oHpfKvgrv9dVIMTVAbY2IL3U0BRBOvrorBSWhQd3X8CAAD//w4a/Ri9
EQAA
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
		_escData["/index.html"],
		_escData["/main.css"],
		_escData["/main.js"],
	},
}
