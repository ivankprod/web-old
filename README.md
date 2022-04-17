# ivankprod.ru
IVANKPROD.RU website (golang, js, docker)

[![CI](https://github.com/ivankprod/ivankprod.ru/actions/workflows/ci.yml/badge.svg)](https://github.com/ivankprod/ivankprod.ru/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/ivankprod/ivankprod.ru/branch/main/graph/badge.svg?token=NLBM9MA475)](https://codecov.io/gh/ivankprod/ivankprod.ru)
[![Go Report Card](https://goreportcard.com/badge/github.com/ivankprod/ivankprod.ru/src/server)](https://goreportcard.com/report/github.com/ivankprod/ivankprod.ru/src/server)
![Lines of code](https://img.shields.io/tokei/lines/github/ivankprod/ivankprod.ru)

# Setup
1. Setup
    1. Install Git, Docker, Go and NodeJS
    2. Clone this repo and cd into it

2. Install
```shell
./install.sh
```

3. Build and run in Docker
```shell
# development build
./compose.sh dev

# production build
./compose.sh [prod]
```

4. Generate sitemap.xml (make shure website is running)
```shell
./sitemap.sh
```

5. Rebuild app image to catch sitemap.xml
```shell
# development build
./compose.sh dev

# production build
./compose.sh [prod]
```

# TODOS:
1. user profile page
2. user auth cabinet
