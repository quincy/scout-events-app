FROM gradle:8.6-jdk17 AS build
COPY --chown=gradle:gradle . /home/gradle/src
WORKDIR /home/gradle/src
RUN gradle buildFatJar --no-daemon

FROM openjdk:17
EXPOSE 8080:8080
RUN mkdir /app
COPY --from=build /home/gradle/src/build/libs/scout-events-app-all.jar /app/app.jar
RUN curl --create-dirs -o /app/root.crt 'https://cockroachlabs.cloud/clusters/0de3351e-57c1-4910-836d-5504d3dae7fc/cert'
ENTRYPOINT ["java", "-Xmx200m", "-jar", "/app/app.jar"]
