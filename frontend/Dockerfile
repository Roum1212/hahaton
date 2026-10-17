# syntax=docker/dockerfile:1
FROM node:24-bookworm-slim AS build
WORKDIR /app
ENV NUXT_TELEMETRY_DISABLED=1
RUN npm install --global pnpm@10.32.1
COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
RUN --mount=type=cache,id=hahaton-pnpm,target=/pnpm/store \
    pnpm install --frozen-lockfile --ignore-scripts --store-dir=/pnpm/store
COPY . .
ARG APP_PUBLIC_SITE_NAME="Цифровой волонтёр Обнинска"
ARG APP_PUBLIC_SITE_DESCRIPTION="Добрые дела, волонтёрские часы и поощрения"
ENV APP_PUBLIC_SITE_NAME=$APP_PUBLIC_SITE_NAME \
    APP_PUBLIC_SITE_DESCRIPTION=$APP_PUBLIC_SITE_DESCRIPTION \
    NITRO_PRESET=node-server
RUN pnpm rebuild && pnpm build

FROM node:24-bookworm-slim AS runtime
WORKDIR /app
ENV NODE_ENV=production \
    NITRO_HOST=0.0.0.0 \
    NITRO_PORT=3000
COPY --from=build --chown=node:node /app/.output ./.output
USER node
EXPOSE 3000
HEALTHCHECK --interval=10s --timeout=5s --start-period=10s --retries=6 \
    CMD node -e "fetch('http://127.0.0.1:3000/').then(r=>process.exit(r.ok?0:1)).catch(()=>process.exit(1))"
CMD ["node", ".output/server/index.mjs"]
