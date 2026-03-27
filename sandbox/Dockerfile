# syntax=docker/dockerfile:1
FROM debian:bookworm-slim

# Install system dependencies
RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        git \
        curl \
        ca-certificates \
        gnupg \
        python3 \
        nodejs \
        npm \
    && rm -rf /var/lib/apt/lists/*

# Install GitHub CLI
RUN curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg \
        | gpg --dearmor -o /usr/share/keyrings/githubcli-archive-keyring.gpg \
    && echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" \
        > /etc/apt/sources.list.d/github-cli.list \
    && apt-get update \
    && apt-get install -y --no-install-recommends gh \
    && rm -rf /var/lib/apt/lists/*

# Create non-root user
RUN useradd -m -s /bin/bash agent

USER agent
WORKDIR /home/agent

# Install Claude Code via official install script
RUN curl -fsSL https://claude.ai/install.sh | bash
ENV PATH="/home/agent/.local/bin:${PATH}"

# Copy entrypoint scripts
COPY --chown=agent:agent entrypoint-common.sh entrypoint-common.sh
COPY --chown=agent:agent entrypoint-drive.sh entrypoint-drive.sh
COPY --chown=agent:agent entrypoint-coast.sh entrypoint-coast.sh
RUN chmod +x entrypoint-common.sh entrypoint-drive.sh entrypoint-coast.sh
