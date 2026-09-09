# syntax=docker/dockerfile:1

# -------------------------------------------------------------------
# Stage 1: Builder using Bun
# -------------------------------------------------------------------
FROM oven/bun:1-alpine AS builder

WORKDIR /app

# Install dependencies with frozen lockfile for reproducible builds
COPY package.json bun.lock* ./
RUN bun install --frozen-lockfile

# Copy source files
COPY . .

# Build SvelteKit application for production
RUN bun run build

# Prune development dependencies for minimal production footprint
RUN rm -rf node_modules && bun install --production --frozen-lockfile

# -------------------------------------------------------------------
# Stage 2: Lean Production Runner using Bun
# -------------------------------------------------------------------
FROM oven/bun:1-alpine AS runner

WORKDIR /app

# Create unprivileged user for security
RUN adduser -D -g '' -u 1001 bunuser

# Copy built application and production dependencies
COPY --from=builder /app/package.json ./
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/build ./build
COPY --from=builder /app/static ./static

# Set permissions
RUN chown -R bunuser:bunuser /app

USER bunuser:bunuser

ENV NODE_ENV=production
ENV PORT=3000
ENV HOST=0.0.0.0

EXPOSE 3000

CMD ["bun", "run", "./build/index.js"]
