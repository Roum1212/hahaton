import withNuxt from './.nuxt/eslint.config.mjs';

export default withNuxt(
  {
    ignores: ['gen/**', '.output/**'],
  },
  {
    files: ['**/*.ts'],
    rules: {
      'no-console': 'off',
      '@typescript-eslint/no-namespace': 'off',
    },
  },
  {
    files: ['**/*.vue'],
    rules: {
      '@typescript-eslint/no-empty-object-type': 'off',
      'vue/no-multiple-template-root': 'off',
      '@typescript-eslint/no-namespace': 'off',
      '@typescript-eslint/no-duplicate-enum-values': 'off',
      'vue/no-parsing-error': 'off',
      'vue/no-template-shadow': 'off',
      'vue/html-self-closing': [
        'off',
        {
          html: {
            void: 'never',
          },
        },
      ],
    },
  },
);
