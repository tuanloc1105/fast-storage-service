FROM tuanloc/my_ubuntu:oracle

EXPOSE 8080

ADD ./go_app /go_app
ADD ./additional_source_code/ /additional_source_code/
ADD ./static/ /static/
ADD ./templates/ /templates/

CMD ["/go_app"]
