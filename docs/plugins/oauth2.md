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
    scopes:
      - user
      - repo
      - read:org
```

## Permissionamento

Você poderá usar os times do Github para atribuir permissões de futuras features.
