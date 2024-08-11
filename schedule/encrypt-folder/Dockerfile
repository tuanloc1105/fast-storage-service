FROM tuanloc/my_ubuntu:oracle

EXPOSE 8080

COPY ./service_scheduler.py /service_scheduler.py

RUN pip3.12 install schedule
RUN pip3.12 install requests

CMD ["python3.12", "service_scheduler.py"]
