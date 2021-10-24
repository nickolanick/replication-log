FROM python:3

WORKDIR /python
COPY requirements.txt requirements.txt 

RUN pip3 install -r requirements.txt

COPY . /python

ENTRYPOINT python3 main.py
