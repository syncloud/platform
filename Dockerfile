FROM python:3.8-buster
COPY requirements.txt /
RUN apt update
RUN apt install -y wget build-essential libsasl2-dev libldap2-dev libssl-dev libjansson-dev
RUN pip install -r /requirements.txt
RUN rm /requirements.txt
RUN rm -rf /var/lib/apt/lists/*
