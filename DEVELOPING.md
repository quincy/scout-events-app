# How to set up your development environment

## SDKs

This project uses [SDKMAN!](https://sdkman.io/) to manage things like the Gradle version.
See the [.sdkmanrc](.sdkmanrc) file to see which versions of the tools it will install.

To get started with SDKMAN!, visit https://sdkman.io/install and follow the installation instructions.

You can install by running this command on Unix-like operating systems:

```shell
curl -s "https://get.sdkman.io" | bash
```

You can install the tools listed in the [.sdkmanrc](.sdkmanrc) file with:

```shell
sdk env install
```

You can activate these tools with:

```shell
sdk env
```

Or you can set the following configuration in your `~/.sdkman/etc/config` file:

```text
sdkman_auto_env=true
```


## Docker

If you want to run the app locally in a docker container, you'll need docker. Follow the instructions
at [docker.com](https://docs.docker.com/get-docker/).


## Build

If you're using Intellij, there are Run Configurations available in the repository:

- [pre-commit build](.idea/runConfigurations/pre_commit_build.xml)
  - run before committing
  - ```shell
    ./gradlew clean spotlessApply detekt build check
    ```
- [build docker image](.idea/runConfigurations/build_docker_image.xml)
  - run to build the docker file locally
  - ```shell
    docker build -t scout-events-app .
    ```


## Running the app locally

You'll need to first create the docker-compose environment, which will start a local version of the database:

```shell
docker-compose up
```

Use the [run app](.idea/runConfigurations/run_app.xml) Run Configuration in IntelliJ.

If you want to run the docker container you built above:

```shell
docker run \
  --cpus=1 --memory=512m \
  -p 8080:8080 \
  -e DB_HOSTNAME=scout-events-db \
  --network=scout-events-app_scout-events-net \
  scout-events-app
```


### Running liquibase

Liquibase is used to manage changes to the database schema.  Liquibase is run automatically in the docker compose
environment.  But you can run it manually as well.

Check the status:

```shell
liquibase status --username=liquibase --url=jdbc:postgresql://localhost:26257/scouting --changelog-file=liquibase/changelog/changelog-root.yaml
```

Check what SQL will be run:

```shell
liquibase update-sql --username=liquibase --url=jdbc:postgresql://localhost:26257/scouting --changelog-file=liquibase/changelog/changelog-root.yaml
```

Run the liquibase migration:

```shell
liquibase update --username=liquibase --url=jdbc:postgresql://localhost:26257/scouting --changelog-file=liquibase/changelog/changelog-root.yaml
```
