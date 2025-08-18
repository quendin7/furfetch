pkgname=asfetch
pkgver=2
pkgrel=1
pkgdesc="A fast system information fetching tool written in Go, similar to Neofetch."
arch=('x86_64')
url="https://github.com/quendin7/asf"
license=('MIT')
depends=('git')
makedepends=('go')
optdepends=(
  'playerctl: for Spotify/music info'
  'lsb-release: for detailed OS info'
  'pciutils: for GPU info (lspci)'
  'upower: for battery info'
  'vulkan-tools: for detailed GPU info'
)
source=("git+${url}.git#tag=v${pkgver}?subdir=${pkgname}")
sha256sums=('SKIP')

build() {
    cd "${srcdir}/${pkgname}"
    go mod tidy
    CGO_ENABLED=0 go build -ldflags="-s -w" -o asf
}

package() {
    cd "${srcdir}/${pkgname}"
    install -D -m755 asf "${pkgdir}/usr/bin/asfetch"
}
