FROM alpine:3.13
RUN addgroup gimmesubnet && adduser -S -G gimmesubnet gimmesubnet
COPY gimmesubnet /bin/gimmesubnet
USER gimmesubnet
ENTRYPOINT ["/bin/gimmesubnet"]
