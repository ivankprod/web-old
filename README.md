# ivankprod.ru
IVANKPROD.RU website
(go, webpack)

[![CI](https://github.com/ivankprod/ivankprod.ru/actions/workflows/ci.yml/badge.svg)](https://github.com/ivankprod/ivankprod.ru/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/ivankprod/ivankprod.ru/branch/main/graph/badge.svg?token=NLBM9MA475)](https://codecov.io/gh/ivankprod/ivankprod.ru)

# Setup
1. Setup
    1. Install Git, Docker, Go and NodeJS
    2. Clone this repo and cd into it

2. Install
```shell
./install.sh
```
3. Build to /build_(dev|prod) dir
```shell
# development
./build.sh dev [os] [arch]

# production
./build.sh [prod] [os] [arch]
```

4. Run in Docker
```shell
# development build
./compose.sh dev

# production build
./compose.sh [prod]
```

5. Generate sitemap.xml for (dev|prod) build (make shure website is running)
```shell
# development build
./sitemap.sh dev

# production build
./sitemap.sh [prod]
```

# TODOS:
1. user profile page
2. user auth cabinet

