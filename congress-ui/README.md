# React + TypeScript + Vite

This template provides a minimal setup to get React working in Vite with HMR and some ESLint rules.

Currently, two official plugins are available:


## React Compiler

The React Compiler is not enabled on this template because of its impact on dev & build performances. To add it, see [this documentation](https://react.dev/learn/react-compiler/installation).

## Expanding the ESLint configuration

If you are developing a production application, we recommend updating the configuration to enable type-aware lint rules:

```js
export default defineConfig([
  globalIgnores(['dist']),
  {
    files: ['**/*.{ts,tsx}'],
    extends: [
      // Other configs...

      // Remove tseslint.configs.recommended and replace with this
      tseslint.configs.recommendedTypeChecked,
      // Alternatively, use this for stricter rules
      tseslint.configs.strictTypeChecked,
      // Optionally, add this for stylistic rules
      tseslint.configs.stylisticTypeChecked,

      // Other configs...
    ],
    languageOptions: {
      parserOptions: {
        project: ['./tsconfig.node.json', './tsconfig.app.json'],
        tsconfigRootDir: import.meta.dirname,
      },
      // other options...
    },
  },
])

```

You can also install [eslint-plugin-react-x](https://npmx.dev/package/eslint-plugin-react-x) and [eslint-plugin-react-dom](https://npmx.dev/package/eslint-plugin-react-dom) for React-specific lint rules:

```js
// eslint.config.js
import reactX from 'eslint-plugin-react-x'
import reactDom from 'eslint-plugin-react-dom'

export default defineConfig([
  globalIgnores(['dist']),
  {
    files: ['**/*.{ts,tsx}'],
    extends: [
      // Other configs...
      // Enable lint rules for React
      reactX.configs['recommended-typescript'],
      // Enable lint rules for React DOM
      reactDom.configs.recommended,
    ],
    languageOptions: {
      parserOptions: {
        project: ['./tsconfig.node.json', './tsconfig.app.json'],
        tsconfigRootDir: import.meta.dirname,
      },
      // other options...
    },
  },
])

```
# Congress UI

Frontend React + TypeScript do MS Congress para contas de pessoa jurídica (`tipoPessoaID = 2`).

## Executar

Com a API rodando em `http://localhost:8080`:

```bash
npm install
npm run dev
```

Abra `http://localhost:5173`.

## Executar com Docker Compose

Na pasta `ms-congress`:

```bash
docker compose up --build
```

Abra `http://localhost:5173`. A API continuará disponível em `http://localhost:8080`.

## Fluxo implementado

- Login em `POST /v1/api/login`.
- Busca do perfil em `GET /v1/api/usuarios/{id}` para obter o ID da pessoa jurídica.
- Listagem em `GET /v1/api/eventos/organizador/{organizadorId}`.
- Detalhes em `GET /v1/api/eventos/{id}`.
- Criação, edição e exclusão de eventos.
- Perfil da organização e logout.
- Sessão persistida no `localStorage`.

## Estrutura

- `src/App.tsx`: composição do painel e estado de navegação.
- `src/services/api.ts`: cliente HTTP e sessão.
- `src/types.ts`: contratos TypeScript da API.
- `src/utils.ts`: formatação de datas e erros.
- `src/components/LoginScreen.tsx`: autenticação.
- `src/components/EventCard.tsx`: item da lista de eventos.
- `src/components/EventForm.tsx`: criação e edição.
