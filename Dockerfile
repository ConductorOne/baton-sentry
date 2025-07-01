FROM gcr.io/distroless/static-debian11:nonroot
ENTRYPOINT ["/baton-sentry"]
COPY baton-sentry /