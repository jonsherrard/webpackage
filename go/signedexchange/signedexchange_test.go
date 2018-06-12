package signedexchange_test

import (
	"bytes"
	"encoding/pem"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	. "github.com/WICG/webpackage/go/signedexchange"
	"github.com/WICG/webpackage/go/signedexchange/internal/testhelper"
)

const (
	payload  = `Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.`
	pemCerts = `-----BEGIN CERTIFICATE-----
MIIF8jCCBNqgAwIBAgIQDmTF+8I2reFLFyrrQceMsDANBgkqhkiG9w0BAQsFADBw
MQswCQYDVQQGEwJVUzEVMBMGA1UEChMMRGlnaUNlcnQgSW5jMRkwFwYDVQQLExB3
d3cuZGlnaWNlcnQuY29tMS8wLQYDVQQDEyZEaWdpQ2VydCBTSEEyIEhpZ2ggQXNz
dXJhbmNlIFNlcnZlciBDQTAeFw0xNTExMDMwMDAwMDBaFw0xODExMjgxMjAwMDBa
MIGlMQswCQYDVQQGEwJVUzETMBEGA1UECBMKQ2FsaWZvcm5pYTEUMBIGA1UEBxML
TG9zIEFuZ2VsZXMxPDA6BgNVBAoTM0ludGVybmV0IENvcnBvcmF0aW9uIGZvciBB
c3NpZ25lZCBOYW1lcyBhbmQgTnVtYmVyczETMBEGA1UECxMKVGVjaG5vbG9neTEY
MBYGA1UEAxMPd3d3LmV4YW1wbGUub3JnMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8A
MIIBCgKCAQEAs0CWL2FjPiXBl61lRfvvE0KzLJmG9LWAC3bcBjgsH6NiVVo2dt6u
Xfzi5bTm7F3K7srfUBYkLO78mraM9qizrHoIeyofrV/n+pZZJauQsPjCPxMEJnRo
D8Z4KpWKX0LyDu1SputoI4nlQ/htEhtiQnuoBfNZxF7WxcxGwEsZuS1KcXIkHl5V
RJOreKFHTaXcB1qcZ/QRaBIv0yhxvK1yBTwWddT4cli6GfHcCe3xGMaSL328Fgs3
jYrvG29PueB6VJi/tbbPu6qTfwp/H1brqdjh29U52Bhb0fJkM9DWxCP/Cattcc7a
z8EXnCO+LK8vkhw/kAiJWPKx4RBvgy73nwIDAQABo4ICUDCCAkwwHwYDVR0jBBgw
FoAUUWj/kK8CB3U8zNllZGKiErhZcjswHQYDVR0OBBYEFKZPYB4fLdHn8SOgKpUW
5Oia6m5IMIGBBgNVHREEejB4gg93d3cuZXhhbXBsZS5vcmeCC2V4YW1wbGUuY29t
ggtleGFtcGxlLmVkdYILZXhhbXBsZS5uZXSCC2V4YW1wbGUub3Jngg93d3cuZXhh
bXBsZS5jb22CD3d3dy5leGFtcGxlLmVkdYIPd3d3LmV4YW1wbGUubmV0MA4GA1Ud
DwEB/wQEAwIFoDAdBgNVHSUEFjAUBggrBgEFBQcDAQYIKwYBBQUHAwIwdQYDVR0f
BG4wbDA0oDKgMIYuaHR0cDovL2NybDMuZGlnaWNlcnQuY29tL3NoYTItaGEtc2Vy
dmVyLWc0LmNybDA0oDKgMIYuaHR0cDovL2NybDQuZGlnaWNlcnQuY29tL3NoYTIt
aGEtc2VydmVyLWc0LmNybDBMBgNVHSAERTBDMDcGCWCGSAGG/WwBATAqMCgGCCsG
AQUFBwIBFhxodHRwczovL3d3dy5kaWdpY2VydC5jb20vQ1BTMAgGBmeBDAECAjCB
gwYIKwYBBQUHAQEEdzB1MCQGCCsGAQUFBzABhhhodHRwOi8vb2NzcC5kaWdpY2Vy
dC5jb20wTQYIKwYBBQUHMAKGQWh0dHA6Ly9jYWNlcnRzLmRpZ2ljZXJ0LmNvbS9E
aWdpQ2VydFNIQTJIaWdoQXNzdXJhbmNlU2VydmVyQ0EuY3J0MAwGA1UdEwEB/wQC
MAAwDQYJKoZIhvcNAQELBQADggEBAISomhGn2L0LJn5SJHuyVZ3qMIlRCIdvqe0Q
6ls+C8ctRwRO3UU3x8q8OH+2ahxlQmpzdC5al4XQzJLiLjiJ2Q1p+hub8MFiMmVP
PZjb2tZm2ipWVuMRM+zgpRVM6nVJ9F3vFfUSHOb4/JsEIUvPY+d8/Krc+kPQwLvy
ieqRbcuFjmqfyPmUv1U9QoI4TQikpw7TZU0zYZANP4C/gj4Ry48/znmUaRvy2kvI
l7gRQ21qJTK5suoiYoYNo3J9T+pXPGU7Lydz/HwW+w0DpArtAaukI8aNX4ohFUKS
wDSiIIWIWJiJGbEeIO0TIFwEVWTOnbNl/faPXpk5IRXicapqiII=
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIIEsTCCA5mgAwIBAgIQBOHnpNxc8vNtwCtCuF0VnzANBgkqhkiG9w0BAQsFADBs
MQswCQYDVQQGEwJVUzEVMBMGA1UEChMMRGlnaUNlcnQgSW5jMRkwFwYDVQQLExB3
d3cuZGlnaWNlcnQuY29tMSswKQYDVQQDEyJEaWdpQ2VydCBIaWdoIEFzc3VyYW5j
ZSBFViBSb290IENBMB4XDTEzMTAyMjEyMDAwMFoXDTI4MTAyMjEyMDAwMFowcDEL
MAkGA1UEBhMCVVMxFTATBgNVBAoTDERpZ2lDZXJ0IEluYzEZMBcGA1UECxMQd3d3
LmRpZ2ljZXJ0LmNvbTEvMC0GA1UEAxMmRGlnaUNlcnQgU0hBMiBIaWdoIEFzc3Vy
YW5jZSBTZXJ2ZXIgQ0EwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQC2
4C/CJAbIbQRf1+8KZAayfSImZRauQkCbztyfn3YHPsMwVYcZuU+UDlqUH1VWtMIC
Kq/QmO4LQNfE0DtyyBSe75CxEamu0si4QzrZCwvV1ZX1QK/IHe1NnF9Xt4ZQaJn1
itrSxwUfqJfJ3KSxgoQtxq2lnMcZgqaFD15EWCo3j/018QsIJzJa9buLnqS9UdAn
4t07QjOjBSjEuyjMmqwrIw14xnvmXnG3Sj4I+4G3FhahnSMSTeXXkgisdaScus0X
sh5ENWV/UyU50RwKmmMbGZJ0aAo3wsJSSMs5WqK24V3B3aAguCGikyZvFEohQcft
bZvySC/zA/WiaJJTL17jAgMBAAGjggFJMIIBRTASBgNVHRMBAf8ECDAGAQH/AgEA
MA4GA1UdDwEB/wQEAwIBhjAdBgNVHSUEFjAUBggrBgEFBQcDAQYIKwYBBQUHAwIw
NAYIKwYBBQUHAQEEKDAmMCQGCCsGAQUFBzABhhhodHRwOi8vb2NzcC5kaWdpY2Vy
dC5jb20wSwYDVR0fBEQwQjBAoD6gPIY6aHR0cDovL2NybDQuZGlnaWNlcnQuY29t
L0RpZ2lDZXJ0SGlnaEFzc3VyYW5jZUVWUm9vdENBLmNybDA9BgNVHSAENjA0MDIG
BFUdIAAwKjAoBggrBgEFBQcCARYcaHR0cHM6Ly93d3cuZGlnaWNlcnQuY29tL0NQ
UzAdBgNVHQ4EFgQUUWj/kK8CB3U8zNllZGKiErhZcjswHwYDVR0jBBgwFoAUsT7D
aQP4v0cB1JgmGggC72NkK8MwDQYJKoZIhvcNAQELBQADggEBABiKlYkD5m3fXPwd
aOpKj4PWUS+Na0QWnqxj9dJubISZi6qBcYRb7TROsLd5kinMLYBq8I4g4Xmk/gNH
E+r1hspZcX30BJZr01lYPf7TMSVcGDiEo+afgv2MW5gxTs14nhr9hctJqvIni5ly
/D6q1UEL2tU2ob8cbkdJf17ZSHwD2f2LSaCYJkJA69aSEaRkCldUxPUd1gJea6zu
xICaEnL6VpPX/78whQYwvwt/Tv9XBZ0k7YXDK/umdaisLRbvfXknsuvCnQsH6qqF
0wGjIChBWUMo0oHjqvbsezt3tkBigAVBRQHvFwY+3sAzm2fTYS5yh+Rp/BIAV0Ae
cPUeybQ=
-----END CERTIFICATE-----
`
	// Generated by `openssl genrsa -out privatekey.pem 2048`
	pemPrivateKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEoAIBAAKCAQEAoMRYVlgUxlVOvejxDblbIZAg4ZtTbAmI7/YzNqmlKBB7UGik
7t6MCTJRM1PAQoDdRC0H5XI0TS04Lizwet8gEeBMtyHqLcWmOUGYNsYO7nNgT7N2
wbEs6v6KHHPHPMKzmxMPayOWrfE7mRvHvwTtIbE5ar5PNjpypjNH24TddkAmIXbM
YbkS2F43rVgpzOihjbeTQ/A6pxqcplifmoGSI6W26dg5N9yGnmo1ZcLdpHixR9Lr
e3xvunkDxT+B0OlwBRQtTQvZ1YoDWylpq3cOiFqU0Wn9+AG8JpL2yI49KQMVKyBV
7dLtr43LFhtBefkyqSNTxqPZyUAJJ2SNkJgwIQIDAQABAoIBAFJz4QqHqj/+SKBF
9DuhsQeJsBOFYkeqrDzF/IYwg7AEo/odcVnBcfjVgafdcGGrTdBFeCNJa2GZq5Kj
IcMi5IPGkhHqpvxKvnHnHnYZJldNfTvjQykcAXmUiqkFCE41XYBPSj0cx472hiaE
hPGHSUdaaaRBbsbVOy/aZSRFBIA8ngxyrW6B94Q/uLVZBn6axqoj8xT1YFVBgH5G
/lVxfkpjUD2im9r3w+7ofSmMKa6CyJ/bBdRf8p0ACyzDbkfyXjwUxSj/ZFrpLg66
amEXgauqxKEAhF8MP8oKir9aEwl7EaYFIRFpzQ6LT6edD5vcieov6hDi1f8xxdty
5lL4HkECgYEA01+pVvn2VqANu9tgpcX3srY6QKnqViBSXr6GX+XpcCJlxR2S4FVD
gdEwMHJK9137krvzIek57BFQXd4bTpeUW3Da8rX73tUnqKrQ5pmEqpghRyCqo0kT
V1ObepNUcQVmK6VnqIuckHNV7sjYnSCgY4P4WiPBRJCG3jTI2LUpo/UCgYEAwrV9
MtwsV9HlVHNrd8hqqaXnDvY1InFCfFxyR0m5KMTiwvcswBbwpTYtKZXWnz2HRVbO
aMmh2RQKk9Swpwb/q2TjVnPPUqH14++OwyR0k/0L4KBZMY736GqyWnfod6G5KQD2
f5MtwRFCYoJ6Tts4KtMzxxaV4TeRQA0EES7rK/0CgYBVztbi7TSYs/7/TS6t/XDx
xtJdH912u0ZVGglY8u/SStR/seLHWTW/hJmIgU13oFqZld083f5anCjBAoKZZCWg
/W6U61XlfyjLaxTFGHtn+bxAsL007lyArftHRnoYK7XvcAVlwc98QKYY+sYc+3rB
C3kNtsglunpVyJ3kg5705QJ/cVMwi2maZYLE92I2KoF7k0H8ObkTM/i3uaoU2WkP
W6s8UD2MzkCLz5y4rHuJbyVglfrwKA0zJiWEAobISm7IX/lYV/kPsgiSFRhY/zs4
numpABRT1YRgxeVT6VPg+cAnBLaKwbXn63cgLDXE+iCdkE9c04NRuMOexqjMtTOZ
rQKBgDSCTKwnbJUqN94WdBYjinFN/bR6E0wW640jkB/3e8Y4a+W4OVHWlxoEu4Tm
s5B6gZsV/ojttR+aaeRknfrhQwEIA/k2r2oZE9yp8djzyiiqGswgw8yO0WSJztbx
GRqzPwjon7ESIVpKLrVuh5qlMhUkOFUeF9wvViWX4qnV5Fvg
-----END RSA PRIVATE KEY-----
`
	expectedSignatureHeader = "label; sig=*WBpXmSAfHvg0ONKuqGy1sI0+U0j7kE7ZoclSxr1/VZwvHvwU5exZR2jD7HnCMkvSQBtoMseXkXwy/75xwDfksIH+vpGPyb51FNwONOO9HwCzlIN3IM1KIqgm/OQNKoNJs7CGZw9S+m+aj+n5cq6v30SU46jlgEn1qXzEhPLIE2x2eK2rBkk76ifqgjUuAwY46jOwq3ihjdGBakvbTrxjhGbhI3O1bn3Ueaucx5/dU+UKl+XqFq1r0kxIUF2KLXckPVYG4hVj2eHS9sGjG8C1vEDboxgcjB9lcmoryDCTCKfixQwSoZ+c84sb1x2rjjzuf8NskKReTsUExsiGsh9/Yg==*; validity-url=\"https://example.com/resource.validity\"; integrity=\"mi\"; cert-url=\"https://example.com/cert.msg\"; cert-sha256=*ZC3lTYTDBJQVf1P2V7+fibTqbIsWNR/X7CWNVW+CEEA=*; date=1517418800; expires=1517422400"
)

type zeroReader struct{}

func (zeroReader) Read(b []byte) (int, error) {
	for i := range b {
		b[i] = 0
	}
	return len(b), nil
}

func mustReadFile(path string) []byte {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return b
}

func TestSignedExchange(t *testing.T) {
	u, _ := url.Parse("https://example.com/")
	header := http.Header{}
	header.Add("Content-Type", "text/html; charset=utf-8")

	// Multiple values for the same header
	header.Add("Foo", "Bar")
	header.Add("Foo", "Baz")

	e, err := NewExchange(u, nil, 200, header, []byte(payload))
	if err != nil {
		t.Fatal(err)
	}
	if err := e.MiEncodePayload(16); err != nil {
		t.Fatal(err)
	}

	now := time.Date(2018, 1, 31, 17, 13, 20, 0, time.UTC)
	certs, err := ParseCertificates([]byte(pemCerts))
	if err != nil {
		t.Fatal(err)
	}

	derPrivateKey, _ := pem.Decode([]byte(pemPrivateKey))
	privKey, err := ParsePrivateKey(derPrivateKey.Bytes)
	if err != nil {
		t.Fatal(err)
	}
	certUrl, _ := url.Parse("https://example.com/cert.msg")
	validityUrl, _ := url.Parse("https://example.com/resource.validity")
	s := &Signer{
		Date:        now,
		Expires:     now.Add(1 * time.Hour),
		Certs:       certs,
		CertUrl:     certUrl,
		ValidityUrl: validityUrl,
		PrivKey:     privKey,
		Rand:        zeroReader{},
	}
	if err := e.AddSignatureHeader(s); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	if err := WriteExchangeFile(&buf, e); err != nil {
		t.Fatal(err)
	}

	magic, err := buf.ReadBytes(0x00)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(magic, HeaderMagicBytes) {
		t.Errorf("unexpected magic: %q", magic)
	}

	var encodedSigLength [3]byte
	if _, err := io.ReadFull(&buf, encodedSigLength[:]); err != nil {
		t.Fatal(err)
	}
	sigLength := Decode3BytesBigEndianUint(encodedSigLength)

	if sigLength != len(expectedSignatureHeader) {
		t.Errorf("Unexpected sigLength: %d", sigLength)
	}

	var encodedHeaderLength [3]byte
	if _, err := io.ReadFull(&buf, encodedHeaderLength[:]); err != nil {
		t.Fatal(err)
	}
	headerLength := Decode3BytesBigEndianUint(encodedHeaderLength)

	signatureHeaderBytes := make([]byte, sigLength)
	if _, err := io.ReadFull(&buf, signatureHeaderBytes); err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(signatureHeaderBytes, []byte(expectedSignatureHeader)) {
		t.Errorf("Unexpected signature header: %q", signatureHeaderBytes)
	}

	encodedHeader := make([]byte, headerLength)
	if _, err := io.ReadFull(&buf, encodedHeader); err != nil {
		t.Fatal(err)
	}

	got, err := testhelper.CborBinaryToReadableString(encodedHeader)
	if err != nil {
		t.Fatal(err)
	}
	want := strings.TrimSpace(string(mustReadFile("test-signedexchange-expected.txt")))

	if got != want {
		t.Errorf("WriteExchangeFile:\ngot: %v\nwant: %v", got, want)
	}

	gotPayload := buf.Bytes()
	wantPayload := mustReadFile("test-signedexchange-expected-payload-mi.bin")
	if !bytes.Equal(gotPayload, wantPayload) {
		t.Errorf("payload mismatch")
	}
}

func TestSignedExchangeStatefulHeader(t *testing.T) {
	u, _ := url.Parse("https://example.com/")
	header := http.Header{}
	header.Add("Content-Type", "text/html; charset=utf-8")
	// Set-Cookie is a stateful header and not available.
	header.Add("Set-Cookie", "wow, such cookie")

	if _, err := NewExchange(u, nil, 200, header, []byte(payload)); err == nil {
		t.Fatal(err)
	}

	// Header names are case-insensitive.
	u, _ = url.Parse("https://example.com/")
	header = http.Header{}
	header.Add("cOnTent-TyPe", "text/html; charset=utf-8")
	header.Add("setProfile", "profile X")

	if _, err := NewExchange(u, nil, 200, header, []byte(payload)); err == nil {
		t.Fatal(err)
	}
}
