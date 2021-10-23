# ivankprod.ru
IVANKPROD.RU website
(go, webpack)

# Setup
1. Setup
  1. Install Git, Docker, golang and NodeJS
  2. Clone this repo and cd into it
2. Install
```shell
./install.sh
```
3. Build to /build_(dev|prod) dir
```shell
# development
./build.sh dev

# production
./build.sh [prod] [os] [arch]
```

4. Generate sitemap.xml for (dev|prod) build
```shell
# development
./sitemap.sh dev

# production
./sitemap.sh [prod]
```

5. Run in Docker
```shell
# development build
./compose.sh dev

# production build
./compose.sh [prod]
```

# TODOS:
1. user profile page
2. user auth cabinet
