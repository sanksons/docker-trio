FROM ubuntu:16.04
MAINTAINER glitter.sankalp@gmail.com

RUN apt-get update && apt-get install -y golang nodejs supervisor curl htop vim
RUN apt-get install -y mongodb
RUN mkdir -p /env /scripts /var/log/supervisor /var/log/styloko /styloko /var/log/org /org

COPY supervisor.conf /etc/supervisor/conf.d/supervisord.conf
COPY services/styloko/styloko-code /styloko
COPY services/styloko/start-styloko /scripts/start-styloko
COPY services/styloko/styloko.env /env/styloko

COPY services/org/org-code /org
COPY services/org/start-org /scripts/start-org
COPY services/org/org.env /env/org

CMD ["/usr/bin/supervisord"]

