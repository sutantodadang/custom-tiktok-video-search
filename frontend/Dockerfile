FROM node:20-alpine AS base
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable
COPY . /app
WORKDIR /app

FROM base AS build
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --frozen-lockfile
RUN pnpm run build


FROM nginx AS production-stage
EXPOSE 8000
RUN mkdir /app

COPY --from=build /app/dist /app
COPY nginx.conf /etc/nginx/nginx.conf

