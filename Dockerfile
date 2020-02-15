FROM    python:3-slim
RUN     adduser --disabled-password --gecos '' twitch-bot
COPY    . /opt/bot
WORKDIR /opt/bot
RUN     python -m pip install -r requirements.txt
USER    twitch-bot
CMD     ["python", "run.py"]