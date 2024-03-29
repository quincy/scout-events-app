FROM gradle:8.6-jdk17 AS build
COPY --chown=gradle:gradle . /home/gradle/src
WORKDIR /home/gradle/src
RUN gradle buildFatJar --no-daemon

FROM openjdk:17
EXPOSE 8080:8080
RUN mkdir /app
COPY --from=build /home/gradle/src/build/libs/scout-events-app-all.jar /app/app.jar
ENTRYPOINT ["java", "-Xmx200m", "-jar", "/app/app.jar"]
