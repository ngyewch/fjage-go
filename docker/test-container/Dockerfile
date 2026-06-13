FROM eclipse-temurin:11

WORKDIR /workspace/test-container

COPY ./gradlew ./gradlew
COPY ./gradlew.bat ./gradlew.bat
COPY ./gradle ./gradle
RUN ./gradlew --no-daemon

COPY ./build.gradle.kts ./build.gradle.kts
COPY ./settings.gradle.kts ./settings.gradle.kts
RUN ./gradlew dependencies --no-daemon

COPY . .
RUN ./gradlew install --no-daemon

FROM eclipse-temurin:11

COPY --from=0 /workspace/test-container/build/install/test-container /opt/test-container
RUN ls -lR /opt

WORKDIR /workspace/test-container

ENTRYPOINT ["/opt/test-container/bin/test-container"]
