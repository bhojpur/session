FROM moby/buildkit:v0.9.3
WORKDIR /session
COPY session README.md /session/
ENV PATH=/session:$PATH
ENTRYPOINT [ "/bhojpur/session" ]