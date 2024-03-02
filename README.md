# scout-events-app


## Dev Setup

### SDKMAN!
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

You can activate these tools by with:
```shell
sdk env
```

Or you can set the following configuration in your `~/.sdkman/etc/config` file:
```text
sdkman_auto_env=true
```

## Build
To build the project use the following gradle command:
```shell
./gradlew clean spotlessApply detekt build check
```

If you're using IntellIj, a Run Configuration called 
[pre-commit build](.idea/runConfigurations/pre_commit_build.xml) is included.
