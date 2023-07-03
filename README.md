# www

1. run
```
yarn create react-app opent1d --template typescript
mv opent1d www
```

1. remove "react-app" from package.json:
```
  "eslintConfig": {
    "extends": [
      "react-app",
      "react-app/jest"
    ]
  },
```

1. run
```
yarn add @types/testing-library__jest-dom eslint-config-react-app -D
```

1. install sdk
```
yarn dlx @yarnpkg/sdks vscode
```