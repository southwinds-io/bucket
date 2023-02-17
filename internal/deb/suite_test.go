package deb

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"testing"
)

func TestMarshalConfig(t *testing.T) {
	c := &Config{
		Repositories: []Repository{
			{
				Name:         "artisan-enterprise",
				Distribution: "all",
				KeyRef:       "SOUTHWINDS",
			},
		},
		Keys: []Key{
			{
				Ref: "SOUTHWINDS",
				Public: `-----BEGIN PGP PRIVATE KEY BLOCK-----
Version: artisan deb cli-extension: 1.0 build 20230215160800868-65c2b4a47b
Comment: builder key

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
b+mRMZZM9SDySN41etWu0yRm
=dypb
-----END PGP PRIVATE KEY BLOCK-----
`,
				Private: `-----BEGIN PGP PUBLIC KEY BLOCK-----
Version: artisan deb cli-extension: 1.0 build 20230215160800868-65c2b4a47b
Comment: builder key

xsFNBGPuMq0BEAC0i4+2FoGtT+sVy0XBqqgLo1ksXzvIjTQtL3g4QHK/NJnPq9FC
/sc3AnFDhRMcYF7RPS/ikZGN2qQ3tlsEM/VwUvKZC8oVrl4L5EFiONmfmZMPLR6V
RrdhzP/rGOrp2+TIUQvO7JThYqKJqmoRLPAoU1Ty3KikfUOIJp9BpVYY0J49ZErz
M9OJMit8thewYNNuFvmsYja1A0RooBmxrc2u6J8qoBuxCfKhvcUNI3LC1nKm70u8
S1Xk/+e2ljeMWWyqZA05LAakaZFFfln0kxe5nRKb0hEVSX43nprp/6cEcdyc0vqM
7/hqXaHtebJbLj1GYBo8E3cEb4ni6DqRSE8YQ0yeC/vS6YaAIHFxWQLt3R3Mqg0T
RBEL0a8m/53QCLQEIkU3KfULdaVi93jHwanIyfnXyx32EFG9lrXoRMXEzZwDDGrb
tk7U9OLGoxxGXbAFAftZp20hPPqKpwBCzscVkZrc4Apy1CVMmd+zSO0/0UV5cR9X
3PyIhMpzWxoK+pYqYGmi3APlA+cBRq3RFgPjEuDqnYUb9RjhMmGEhNW1CUIcklCK
mIZTPDJRq8Xr34B8vQzEkwEgt7lr89codYmoZyaRC69V0yUcgrGxPoVUF8BwNgM3
6Hvzb/CD6Qio1e6zRQN6eOfnv+bNev9ODfmRIRYpHDUYyvwwyzunUrJe7QARAQAB
zRVhY21lIDxhZG1pbkBhY21lLmNvbT7CwY0EEwEIAEEFAmPuMq0JkJpx2u9zr/a0
FiEEkL0DKxh9umiS2+1qmnHa73Ov9rQCGwMCHgECGQEDCwkHAhUIAxYAAgUnCQIH
AgAArPoP+gL8qaTx7qs0yb9U9RzCy2MWecOsTJXETDzYd82I8/t0q6VahExzDMNn
BTGv/dL5RZ7SSsLxnCridm7kz7PIyDZRueYUEMYgItUQMdHi30xlaR/6CfzUYsoP
kfG9Q/2G3f4BtzVaQCRxgs+p9OAVwB9P8Yntko92FnkzOYUuuXDaAgUP/SbR2P8b
VpoPmjFzdHh0J4jxPJS4+4OGqEZY4lSvwPFadRgJT35VvFQxLKWTF0+U+R4SCmlw
tLMe+dIVg015JxZh0F3TDbLogeNDUza2sYXgig01jnSKmcQ83YK7/SAhYovsBREi
h2ayi3Wu5wreYieIoRLqzUwz4i8YiNCsyQz5AxbJu2ilaOFS0f3JC2SwEqVUDiyh
qhlMgJDEyxfBj3QHSjsg9dpMgFRhw6Xk/p9ldWMKs4+nuAkxeS3HPLlm5N2+YTDG
owamMIOJtUIHzgb7M5yOtPny8PJWQz8OzE07xQHz/ZwTvgD/32iUCkckXm9sw9Ao
oQY5dXI7h2PQcJX/Nz4f+KJvMLruoeq0+sy0+CNYYtI3OWQ5bz5wtLM+4deackbH
rrp16ktInREGk0qPF94GvcL0cQ7+s8+rfIeWFhREhlfs7/hiuPC95vijq9aco/kF
o3GzGL7QguN9+FbEfYvTamieryrfqqgc8p85UXBhmWc76tHytxvVzsFNBGPuMq0B
EACdYeY0u6505Gpa2dWBMqn/7O3yH07pB/fDdo7oBcjyH09P/0CPK/Dt03PwrXY5
8qCOUZQbMo4cvRKflOgQWCwv/6nOW+AGT87mCyeeFuHOxQCQ9aJMyhFVv4msEGfI
T5lqf2oqiXicp8nluN1vxzrGhtwBX8RuCxwgKY6sK6K2QIveo+p1xH+jwqXasVep
XSuoUZDCCavbWKIV1xaYcnJ3kGZJ3ZcBYOh/MTdnsP72HRx+Cur8MysSgFm6EY/k
eYjfPJQE1GTPGZpv2zOFaZ4YIu5Yyu6uAUIZ8gH5pQ1ET4x0ZJFWOVE4JC14ZO6j
IYemD2tDwwvuF3DyvkxCmVqQZobSQkmDg7AyZborRz1o8eFbzVB+Opmnul0BMbaI
aJY0R3utEO3gn27SiiJqFRZ1PRM9dUi2wWZTVWYy0UFEO6bG/2nbmSWvPCPHLcC8
W5VWNLoyaB9Z8sFbDGhGyFb7Iv8oG6/TSySX6Hbgp8UZwVooNd2gYvLM78hbmw/L
gps4H/uJEpSNkb0aJ7Bl5PbLLphm/QlcaMDAH44oocUItQGgHKLf+U1K8Zd1AmAa
PHzfB5jWxGnxE1D5JaPfX0xoGY5ENF+XnjB+Cs6mBft6C1iqG+Fkxwqlf24Giw8F
cLByLK1If/8ZcNue8CNAOb81b5srevOeB9d607sAPELWHwARAQABwsF2BBgBCAAq
BQJj7jKtCZCacdrvc6/2tBYhBJC9AysYfbpoktvtappx2u9zr/a0AhsMAAAovA//
YBUeZ6G0nTAN3kTczhrUmu1fSEXCe3mKQZFD2p7b/KN8nZHg0FTY3jB354I192XX
JUwGl8yAy2KPYsCdF1fSr0wPLFiLhPFxAACD08vlW+VpQqBKvlt27VdSHFyUSJh6
ShJi7LrlZiyUdL2K3CZggoeJX1JlxXi97DHq95iF0lsvIzhj4PeWKMJ5bPlpmc8a
h+67uYVP8DdxvIovxzF4Zl01swB8/qYHR5a5HNu6zR4HzjHCjjM2NcS6H5y7joFc
BLrAi4IfbACpxpFVGhVNHYFYAB+WEBcQ6/fkeAUJ8eM/fYZiQixUn/xcRq4a92M4
Ph8d0X8+h/zh0kxXAvkh7TxkVOu4kxDVhZn/hS4Z2VGeHLESJwWEkzZNsV5GBagq
jPtdy9w54wvizPCz5gQHIXfP0r8M9XlD/zP0KbsV0iy5Fy0f4EGxxUFTYZsH6CVu
2xHs/xtHGz8bjcMZGpeRsJGdNRFGBOb+abiNPpyV1lg4tK7CWMnEGnvxCM0oESPJ
hSrQfhuoLByJGAHDas2pcHHnq6Xu/h9WTw+phOB2OMsCEWOa5T+Q1srw7nydn4SM
nWOxlhT5Pe8Lj8EfTT7l83CqTP/dF7QzaaJCDvAtXALlos9Rv2yXsPYD4fqcZFkH
wjoUufdE7I1+gn5z195v6ZExlkz1IPJI3jV61a7TJGY=
=dFUZ
-----END PGP PUBLIC KEY BLOCK-----
`,
				Passcode: "EDnj!LS5Qc6bKG&Blh9e%1RiVPWJNCTs",
			},
		},
	}
	b, err := yaml.Marshal(c)
	if err != nil {
		t.Fatal(err)
	}
	_ = os.WriteFile("config.yaml", b, 0755)
}

func TestParsePackages(t *testing.T) {
	content, err := os.ReadFile("test/Packages")
	if err != nil {
		t.Fatal(err)
	}
	list := newPackagesDataFromContent(string(content[:]))
	fmt.Println(list)
	list.SaveGz("test/Packages.gz")
}
