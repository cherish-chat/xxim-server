package xrsa

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecrypt(t *testing.T) {

}

func TestEncrypt(t *testing.T) {
	encrypt, err := Encrypt([]byte("hello"), []byte(`-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAuiJyMvMqTKKC5Z4qWU3v
R9ZWm1JnEEP66xYC0a62XsNE+Vi/3OtKChrhsLGzzEfpAmYLtIdODK/Wm5VQeFqA
w/2UgWtIxPrKfLllA3tTcKkbw/K/9WkO24FKmPPg00L7OaVbfvg/0TorLnMyQ65R
OnG8fvs+LqrIRIDgGZIPGCytV4IdV988v/7KHLNUvyAoINLVIISriUwwr5cjAORL
RLsPVW0jJp4xNleE55Vi+0PlmloPwGtEt9xMRIaTIQzpgBzuLLymxF5a5ifbHg/V
xqDumvu1sYCot9fhDqktYsVz990FgpHJv7xeY11ZFvfKYl4T0VLg5Mvzq8+BX5ut
SQIDAQAB
-----END PUBLIC KEY-----
`))
	if err != nil {
		t.Fatalf("Error while encrypting: %v", err)
	}
	bytes, err := Decrypt(encrypt, []byte(`-----BEGIN PRIVATE KEY-----
MIIEvwIBADANBgkqhkiG9w0BAQEFAASCBKkwggSlAgEAAoIBAQC6InIy8ypMooLl
nipZTe9H1labUmcQQ/rrFgLRrrZew0T5WL/c60oKGuGwsbPMR+kCZgu0h04Mr9ab
lVB4WoDD/ZSBa0jE+sp8uWUDe1NwqRvD8r/1aQ7bgUqY8+DTQvs5pVt++D/ROisu
czJDrlE6cbx++z4uqshEgOAZkg8YLK1Xgh1X3zy//socs1S/ICgg0tUghKuJTDCv
lyMA5EtEuw9VbSMmnjE2V4TnlWL7Q+WaWg/Aa0S33ExEhpMhDOmAHO4svKbEXlrm
J9seD9XGoO6a+7WxgKi31+EOqS1ixXP33QWCkcm/vF5jXVkW98piXhPRUuDky/Or
z4Ffm61JAgMBAAECggEABVaBQGYaH5fMgBgpiwvLXo77LY0IkOr/w+Skc/RAkIly
J0imh4frqzKeQZwDAqSwSuxd9Kr2fJiOI48w/zz9mUFJLdUqDjrj+jCZGnQbtNNa
TSPSlSlFyHRvgEeELWk4sHlGKhpYvpMRhto9krsDYzuJeWAsP+UP0h+BB1kZm9Er
/6lneQOZV47zb4tcV0NatY0VVkjHiys/PG9myBnSz2lM27Mr4dq2JzD0ly7g5GEm
aQ9Uty0TILt4Fw6uUu5MwyJ7KcpPSxHwSJEk7HoTEIkk8kRo2x39neIC5le11r8Y
FG583u8yuCZVuitR2Hn+qqar4JVsUr9p51QnhG5CtQKBgQDxzENB4+EG9drXPRmh
hUx1iHySCi0wP0oI6MizsI5LkhoC3Rr+rSKliUZ99uu/GIgGTXXjneTInkElqomp
Wm4QjO378O56NISrIVau5ghmRMbHw9jGYosmq0PYUmxywNG+c8MoTBWvK4nv6S3d
sntD4JwFEIPgBjLOwwqg0WqJ9QKBgQDFETbs1L10fHC7K8UA14nVrpn8wLfDwodC
WMcQRd4T9gm8stUqVbpyR4AU5zbB1lHEGQ0J7LzAFN3J2/qWJ6c+6X4xiogaYXBx
g2RA7T4/ZzxJZTIo8SwKudVpFvqfO6kabHrU6+MCkLrZXFsbqiOFF0HMYlVP+H5Y
O/cOuBRdhQKBgQCmOtY2MzKdtWHIvWGkkF57IuT5BXQBilEchOSN3CBHRd8J/vWz
BlYeVZyXtqxlyLijFJeqbv89CMci3PYc7mVijXCC1yUr8HUQrS/Jt60ombnK1hJu
eIrPf8h5rFiQH41SkIUna/8wWQ9QVw9ILY7eoEjClpMC7V/6k034N2A2DQKBgQCJ
pYyHv9DVHFZhZiEkhWhxKJPGR5YT1jxDy16/rw1/Q8tpUkAyYc7pI6gC8bz9h4V1
Q0ooNINiZzDDXjOZzfizqMPMNsb6JjU0FGJiN1PTVXh2i4iNsGbi1wqJbNOBhVqI
al0he+IEWLMqP6gjmqNUwvnimIyeXyNg3gGi9lDDQQKBgQC8SDKPFiSVhLxBMAqH
VTJjqxEky+NyeCfZz4IcITOZXO7/Jzyu7JrIVmNwkgOYCa3fjl1znZwTK88cH8G3
qQYwMzBz1yCzMLymRC/3X1j8sGiCzkb9prh1Q11zrODwXyGSWB3xnF/EczUBFFMA
ap9exGxwPbbGlYEFo3dYDkBumA==
-----END PRIVATE KEY-----
`))
	if err != nil {
		t.Fatalf("Error while decrypting: %v", err)
	}
	assert.Equal(t, "hello", string(bytes))
}
