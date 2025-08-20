import js from '@eslint/js';
import globals from 'globals';
import react from 'eslint-plugin-react';
import reactHooks from 'eslint-plugin-react-hooks';
import reactRefresh from 'eslint-plugin-react-refresh';
import tseslint from '@typescript-eslint/eslint-plugin';
import tsParser from '@typescript-eslint/parser';
import prettier from 'eslint-plugin-prettier';
import prettierConfig from 'eslint-config-prettier';

export default [
  // Ignore patterns
  { ignores: ['dist', 'node_modules', 'coverage', '*.min.js'] },

  // TypeScript files configuration
  {
    files: ['**/*.{ts,tsx}'],
    languageOptions: {
      ecmaVersion: 2022,
      sourceType: 'module',
      globals: {
        ...globals.browser,
        ...globals.es2022,
        React: 'readonly',
        JSX: 'readonly',
      },
      parser: tsParser,
      parserOptions: {
        ecmaFeatures: { jsx: true },
        project: './tsconfig.app.json',
        tsconfigRootDir: import.meta.dirname,
      },
    },
    settings: {
      react: { version: '18.3' },
    },
    plugins: {
      react,
      'react-hooks': reactHooks,
      'react-refresh': reactRefresh,
      '@typescript-eslint': tseslint,
      prettier,
    },
    rules: {
      // Base JS rules
      ...js.configs.recommended.rules,

      // React rules
      ...react.configs.recommended.rules,
      ...react.configs['jsx-runtime'].rules,
      ...reactHooks.configs.recommended.rules,

      // Prettier integration
      ...prettierConfig.rules,
      'prettier/prettier': 'error',

      // TypeScript rules
      '@typescript-eslint/no-unused-vars': [
        'error',
        {
          argsIgnorePattern: '^_',
          varsIgnorePattern: '^_',
          caughtErrorsIgnorePattern: '^_',
        },
      ],
      '@typescript-eslint/no-explicit-any': 'warn',

      // React specific
      'react/jsx-no-target-blank': 'error',
      'react/jsx-key': 'error',
      'react/self-closing-comp': 'error',
      'react-refresh/only-export-components': 'off', // Disable for now

      // General code quality
      'no-console': ['warn', { allow: ['warn', 'error'] }],
      'no-debugger': 'error',
      'no-unused-vars': 'off', // Use TypeScript version instead
      'no-undef': 'off', // TypeScript handles this
    },
  },

  // JavaScript files configuration
  {
    files: ['**/*.{js,jsx}'],
    languageOptions: {
      ecmaVersion: 2020,
      globals: {
        ...globals.browser,
        React: 'readonly',
        JSX: 'readonly',
      },
      parserOptions: {
        ecmaVersion: 'latest',
        ecmaFeatures: { jsx: true },
        sourceType: 'module',
      },
    },
    settings: { react: { version: '18.3' } },
    plugins: {
      react,
      'react-hooks': reactHooks,
      'react-refresh': reactRefresh,
      prettier,
    },
    rules: {
      ...js.configs.recommended.rules,
      ...react.configs.recommended.rules,
      ...react.configs['jsx-runtime'].rules,
      ...reactHooks.configs.recommended.rules,
      ...prettierConfig.rules,
      'prettier/prettier': 'error',
      'react/jsx-no-target-blank': 'off',
      'react-refresh/only-export-components': [
        'warn',
        { allowConstantExport: true },
      ],
      'no-console': ['warn', { allow: ['warn', 'error'] }],
    },
  },

  // Test files configuration
  {
    files: [
      '**/*.test.{js,jsx,ts,tsx}',
      '**/__tests__/**/*.{js,jsx,ts,tsx}',
      '**/setupTests.js',
    ],
    languageOptions: {
      globals: {
        ...globals.browser,
        vi: 'readonly',
        it: 'readonly',
        test: 'readonly',
        expect: 'readonly',
        describe: 'readonly',
        beforeEach: 'readonly',
        afterEach: 'readonly',
        beforeAll: 'readonly',
        afterAll: 'readonly',
        jest: 'readonly',
      },
    },
    rules: {
      '@typescript-eslint/no-explicit-any': 'off',
      'no-console': 'off',
    },
  },

  // Config files
  {
    files: ['*.config.{js,ts}', 'vite.config.{js,ts}'],
    languageOptions: {
      globals: {
        ...globals.node,
      },
    },
    rules: {
      'no-console': 'off',
      '@typescript-eslint/no-unused-vars': 'off',
    },
  },
];
