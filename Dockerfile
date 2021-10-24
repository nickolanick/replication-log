FROM python:3

COPY requirements.txt requirements.txt 
WORKDIR /python

RUN pip3 install -r requirements.txt

COPY . /python

ENTRYPOINT python3 main.py
