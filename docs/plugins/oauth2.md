# OAuth2

> Nota:
>
> Plugin ainda em desenvolvimento

No momento temos a autenticação via Github.

## Configurações do Plugin

Em seu arquivo de configuração você poderá adicionar seu provider de autenticação, exemplo:

```yaml
oauth2:
  - provider: github
    clientId: ${GITHUB_CLIENT_ID}
    clientSecret: ${GITHUB_CLIENT_SECRET}
    authURL: 'https://github.com/login/oauth/authorize'
    tokenURL: 'https://github.com/login/oauth/access_token'
    urlLoginSuccess: 'http://localhost:3000'
    orgPermission: 'dash-ops'
    scopes:
      - user
      - repo
      - read:org
```

## Permissionamento

Você pode adicionar uma organização no parametro `orgPermission`, vamos usar ela para validar se o usuario que esta tentando efetuar o login esta na organização, caso contrario ele recebera um 401 e será deslogado.
