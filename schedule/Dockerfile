FROM tuanloc/my_ubuntu:oracle

EXPOSE 8080

COPY ./run_scheduler.sh /run_scheduler.sh
COPY ./service_scheduler.py /service_scheduler.py

RUN pip3.12 install schedule
RUN pip3.12 install requests

CMD ["/run_scheduler.sh"]
