module.exports = {
  apps: [
    {
      name: 'frontend',
      script: './.output/server/index.mjs',
      env: {
        APP_NODE_ENV: 'development',
      },
      env_production: {
        APP_NODE_ENV: 'production',
      },
    },
  ],
};
