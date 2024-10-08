val ktorVersion: String by project
val kotlinVersion: String by project
val logbackVersion: String by project

plugins {
  kotlin("jvm") version "1.9.22"
  id("idea")
  id("io.ktor.plugin") version "2.3.8"
  id("org.jetbrains.kotlin.plugin.serialization") version "1.9.22"
  id("io.gitlab.arturbosch.detekt").version("1.23.3")
  id("com.diffplug.spotless") version "6.25.0"
}

group = "com.troop77eagle"

version = "0.0.1-dev"

application {
  mainClass.set("com.troop77eagle.ApplicationKt")

  val isDevelopment: Boolean = project.ext.has("development")
  applicationDefaultJvmArgs = listOf("-Dio.ktor.development=$isDevelopment")
}

repositories {
  mavenCentral()
  maven { url = uri("https://maven.pkg.jetbrains.space/public/p/ktor/eap") }
}

val testcontainersVersion = "1.19.8"

dependencies {
  // ktor
  implementation("io.ktor:ktor-client-content-negotiation")
  implementation("io.ktor:ktor-serialization-kotlinx-json")
  implementation("io.ktor:ktor-server-auth-jvm")
  implementation("io.ktor:ktor-server-auth-jwt-jvm")
  implementation("io.ktor:ktor-server-auto-head-response-jvm")
  implementation("io.ktor:ktor-server-call-id-jvm")
  implementation("io.ktor:ktor-server-call-logging-jvm")
  implementation("io.ktor:ktor-server-compression-jvm")
  implementation("io.ktor:ktor-server-config-yaml")
  implementation("io.ktor:ktor-server-content-negotiation")
  implementation("io.ktor:ktor-server-core-jvm")
  implementation("io.ktor:ktor-server-cors-jvm")
  implementation("io.ktor:ktor-server-host-common-jvm")
  implementation("io.ktor:ktor-server-netty-jvm")
  implementation("io.ktor:ktor-server-resources")
  implementation("io.ktor:ktor-server-sessions-jvm")
  implementation("io.ktor:ktor-server-status-pages-jvm")
  implementation("io.ktor:ktor-server-webjars-jvm")

  // logging
  implementation("ch.qos.logback:logback-classic:$logbackVersion")
  implementation("io.github.oshai:kotlin-logging-jvm:6.0.3")

  // webjars
  implementation("org.webjars.npm:htmx.org:1.9.12")
  implementation("org.webjars:jquery:3.7.1")

  // database
  implementation("org.postgresql:postgresql:42.7.3")
  implementation("com.zaxxer:HikariCP:5.1.0")
  implementation(platform("org.jdbi:jdbi3-bom:3.45.1"))
  implementation("org.jdbi:jdbi3-core")
  implementation("org.jdbi:jdbi3-kotlin")
  implementation("org.jdbi:jdbi3-kotlin-sqlobject")

  // datetime
  implementation("org.jetbrains.kotlinx:kotlinx-datetime:0.6.0-RC.2")

  // templating
  implementation("com.github.spullara.mustache.java:compiler:0.9.10")

  // test
  testImplementation("io.ktor:ktor-server-tests-jvm")
  testImplementation("io.kotest.extensions:kotest-assertions-ktor-jvm:2.0.0")
  testImplementation("io.kotest:kotest-runner-junit5:5.8.1")
  testImplementation("org.junit.jupiter:junit-jupiter-params")
  testImplementation("org.jetbrains.kotlin:kotlin-test-junit:$kotlinVersion")
  testImplementation("io.mockk:mockk:1.13.10")
  testImplementation("com.natpryce:hamkrest:1.8.0.1")
  testImplementation("com.ninja-squad:DbSetup-kotlin:2.1.0")

  // testcontainers
  testImplementation("org.testcontainers:testcontainers:$testcontainersVersion")
  testImplementation("org.testcontainers:junit-jupiter:$testcontainersVersion")
  testImplementation("io.kotest.extensions:kotest-extensions-testcontainers:2.0.1")
  testImplementation("org.testcontainers:cockroachdb:$testcontainersVersion")
  testImplementation("org.liquibase:liquibase-core:4.27.0")
}

idea { module { isDownloadSources = true } }

spotless {
  kotlin { ktfmt() }
  kotlinGradle { ktfmt() }

  format("web") {
    target(
        "**/*.html",
        "**/*.css",
        "**/*.js",
        "**/*.ts",
        "**/*.md",
        "**/*.json",
        "**/*.yaml",
        "**/*.yml")
    targetExclude(".idea/**", ".gradle/**", "build/**")
    prettier().config(mapOf("tabWidth" to 2))
  }

  format("build-scripts") {
    // define the files to apply `misc` to
    target("*.gradle", ".gitattributes", ".gitignore")

    // define the steps to apply to those files
    trimTrailingWhitespace()
    indentWithSpaces(2)
    endWithNewline()
  }
}

detekt { config.setFrom("$projectDir/config/detekt/detekt.yml") }

ktor {
  docker {
    jreVersion.set(JavaVersion.VERSION_17)

    localImageName.set("scout-events-app")
    imageTag.set("1.0-dev")
  }
}

tasks.withType<Test>().configureEach { useJUnitPlatform() }
