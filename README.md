### THIS PROJECT HAS BEEN MIGRATED TO https://codeberg.org/spagettikod/opent1d

## Build

```
docker build -t registry.spagettikod.se:8443/switchboard:0.1.0 .
docker build -t registry.spagettikod.se:8443/switchboard:0.1.0 --platform linux/amd64 .
docker buildx build --push --progress plain --tag registry.spagettikod.se:8443/opent1d:0.1.0 --platform linux/amd64,linux/arm64 .
```

# www

## Using ChakraUI
```
corepack enable
corepack prepare yarn@stable --activate
yarn create vite opent1d --template react-ts
mv opent1d www
cd www
yarn
yarn dlx @yarnpkg/sdks vscode
yarn add @chakra-ui/react @emotion/react @emotion/styled framer-motion
yarn add @apollo/client graphql
yarn add -D typescript @graphql-codegen/cli @graphql-codegen/client-preset
yarn add -D @graphql-typed-document-node/core
```

add to vite.config.ts :
```
  server: {
    host: "0.0.0.0"
  }
```
