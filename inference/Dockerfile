FROM python:3.12-slim-trixie

COPY --from=ghcr.io/astral-sh/uv:0.8.22 /uv /uvx /bin/

WORKDIR /app

ADD . .

RUN uv sync --locked

CMD ["uv", "run", "main.py"]
