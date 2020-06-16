FROM hayd/deno:debian-1.1.0

RUN apt-get update \
  && apt-get install -y curl ca-certificates --no-install-recommends \
  && curl -LO https://storage.googleapis.com/kubernetes-release/release/`curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt`/bin/linux/amd64/kubectl \
  && chmod +x ./kubectl \
  && mv ./kubectl /usr/local/bin/kubectl \
  && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY ./src/deps.ts /app
RUN deno cache deps.ts

COPY ./src /app
RUN deno cache mod.ts

ENV ORGS="SKT-AIDevOps"
ENV PROJECT="hera"
ENV KEY="f6knfvhhjc57o54vbxvbceuki24zy5uuc2bdf47r656kq7ipbpra"
ENV ENV_NAME="dev"

CMD ["run", "--allow-run", "--allow-net", "--allow-env", "--cached-only", "mod.ts"]