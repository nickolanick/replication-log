FROM python:3

COPY . /python
WORKDIR /python

RUN pip3 install -r requirements.txt

ENTRYPOINT python3 main.py
