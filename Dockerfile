FROM scratch

ADD dsc-amd64 /dsc

ENTRYPOINT ["/dsc"]
