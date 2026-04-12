FROM eclipse-temurin:11

WORKDIR /workspace/test-container

COPY . .
RUN ./gradlew install

FROM eclipse-temurin:11

COPY --from=0 /workspace/test-container/build/install/test-container /opt/test-container
RUN ls -lR /opt

WORKDIR /workspace/test-container

ENTRYPOINT ["/opt/test-container/bin/test-container"]
