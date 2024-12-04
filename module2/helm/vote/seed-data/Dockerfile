# FROM python:3.9-slim
FROM ubuntu:20.04

# add apache bench (ab) tool
RUN apt update \
    && apt install -y software-properties-common \
    && add-apt-repository ppa:deadsnakes/ppa \
    && apt install -y python3.9 \
    && apt install -y --no-install-recommends  apache2-utils  \
    && rm -rf /var/lib/apt/lists/*
# RUN apt-get update \
#     && apt-get install -y --no-install-recommends apache2-utils \
#     && rm -rf /var/lib/apt/lists/*

WORKDIR /seed

COPY . .

# create POST data files with ab friendly formats
RUN python3 make-data.py

CMD /seed/generate-votes.sh