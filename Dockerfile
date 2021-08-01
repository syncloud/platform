FROM python:3.8-buster
COPY requirements.txt /
RUN pip install -f /requirements.txt
RUN rm /requirements.txt
RUN rm -rf /var/lib/apt/lists/*
