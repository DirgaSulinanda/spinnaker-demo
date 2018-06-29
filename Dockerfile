FROM scratch
ADD spinnaker-demo /
ADD web /web
CMD ["/spinnaker-demo"]
