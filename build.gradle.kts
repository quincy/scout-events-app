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

repositories { mavenCentral() }

dependencies {
  // ktor
  implementation("io.ktor:ktor-server-content-negotiation-jvm")
  implementation("io.ktor:ktor-server-core-jvm")
  implementation("io.ktor:ktor-serialization-kotlinx-json-jvm")
  implementation("io.ktor:ktor-server-call-logging-jvm")
  implementation("io.ktor:ktor-server-call-id-jvm")
  implementation("io.ktor:ktor-server-cors-jvm")
  implementation("io.ktor:ktor-server-compression-jvm")
  implementation("io.ktor:ktor-server-host-common-jvm")
  implementation("io.ktor:ktor-server-status-pages-jvm")
  implementation("io.ktor:ktor-server-resources")
  implementation("io.ktor:ktor-server-auto-head-response-jvm")
  implementation("io.ktor:ktor-server-sessions-jvm")
  implementation("io.ktor:ktor-server-auth-jvm")
  implementation("io.ktor:ktor-server-auth-jwt-jvm")
  implementation("io.ktor:ktor-server-netty-jvm")
  implementation("io.ktor:ktor-server-webjars-jvm")
  implementation("org.webjars:jquery:3.7.1")

  // logging
  implementation("ch.qos.logback:logback-classic:$logbackVersion")
  implementation("io.github.oshai:kotlin-logging-jvm:6.0.3")

  // database
  implementation("org.postgresql:postgresql:42.7.3")
  implementation(platform("org.jdbi:jdbi3-bom:3.45.1"))
  implementation("org.jdbi:jdbi3-core")
  implementation("org.jdbi:jdbi3-kotlin")

  // test
  testImplementation("io.ktor:ktor-server-tests-jvm")
  testImplementation("org.jetbrains.kotlin:kotlin-test-junit:$kotlinVersion")
  testImplementation("io.mockk:mockk:1.13.10")
}

idea { module { isDownloadSources = true } }

spotless {
  kotlin { ktfmt() }

  format("misc") {
    // define the files to apply `misc` to
    target("*.gradle", ".gitattributes", ".gitignore")

    // define the steps to apply to those files
    trimTrailingWhitespace()
    indentWithSpaces(2)
    endWithNewline()
  }
}

ktor {
  docker {
    jreVersion.set(JavaVersion.VERSION_17)

    localImageName.set("scout-events-app")
    imageTag.set("1.0-dev")
  }
}
