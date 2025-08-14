# OAuth2

> Nota:
>
> Plugin ainda em desenvolvimento

No momento temos a autenticação via Github.

## Github APP

Primeiro precisamos configurar o token de acesso seguido os seguintes passos:

- Acesse a aba `Settings` da sua organização;
- No menu lateral clique em `OAuth Apps` na seção `Developer settings;
- Registre um novo aplicativo, passando o Nome, URL do APP e URL de callback do login, exemplo:

![Github APP Config](../img/github-config.png)

> Nota:
> Caso você esteja rodando local, precisa referenciar as portas corretas para o front e da API:
>
> ![Github APP Local Config](../img/github-local-config.png)

- Será gerado dois tokens um `Client ID` e `Client Secret`;
- Copie os código e siga as instruções do proximo passo.

## Configurações do Plugin

Em seu arquivo de configuração você poderá adicionar seu provider de autenticação, exemplo:

```yaml
oauth2:
  - provider: github
    clientId: ${GITHUB_CLIENT_ID}
    clientSecret: ${GITHUB_CLIENT_SECRET}
    authURL: 'https://github.com/login/oauth/authorize'
    tokenURL: 'https://github.com/login/oauth/access_token'
    redirectURL: 'http://localhost:8080/api/oauth/redirect'
    urlLoginSuccess: 'http://localhost:5173'
    orgPermission: 'dash-ops'
    scopes:
      - user
      - repo
      - read:org
```

### Parâmetros de Configuração:

- **`clientId`**: Client ID gerado pelo GitHub OAuth App
- **`clientSecret`**: Client Secret gerado pelo GitHub OAuth App
- **`authURL`**: URL de autorização do GitHub
- **`tokenURL`**: URL para troca do código por token
- **`redirectURL`**: URL de callback configurada no GitHub OAuth App (deve apontar para a API)
- **`urlLoginSuccess`**: URL para onde o usuário será redirecionado após login bem-sucedido (frontend)
- **`orgPermission`**: Nome da organização GitHub para validação de permissão
- **`scopes`**: Escopos de permissão solicitados ao GitHub

## Permissionamento

Você pode adicionar uma organização no parametro `orgPermission`, vamos usar ela para validar se o usuario que esta tentando efetuar o login esta na organização, caso contrario ele recebera um 401 e será deslogado.
