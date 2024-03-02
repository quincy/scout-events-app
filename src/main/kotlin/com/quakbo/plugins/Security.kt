package com.quakbo.plugins

import com.auth0.jwt.JWT
import com.auth0.jwt.algorithms.Algorithm
import io.ktor.server.application.Application
import io.ktor.server.application.call
import io.ktor.server.application.install
import io.ktor.server.auth.authentication
import io.ktor.server.auth.jwt.JWTPrincipal
import io.ktor.server.auth.jwt.jwt
import io.ktor.server.response.respondText
import io.ktor.server.routing.get
import io.ktor.server.routing.routing
import io.ktor.server.sessions.Sessions
import io.ktor.server.sessions.cookie
import io.ktor.server.sessions.get
import io.ktor.server.sessions.sessions
import io.ktor.server.sessions.set

fun Application.configureSecurity() {
  data class MySession(val count: Int = 0)
  install(Sessions) { cookie<MySession>("MY_SESSION") { cookie.extensions["SameSite"] = "lax" } }
  // Please read the jwt property from the config file if you are using EngineMain
  val jwtAudience = "jwt-audience"
  val jwtDomain = "https://jwt-provider-domain/"
  val jwtRealm = "ktor sample app"
  val jwtSecret = "secret"
  authentication {
    jwt {
      realm = jwtRealm
      verifier(
          JWT.require(Algorithm.HMAC256(jwtSecret))
              .withAudience(jwtAudience)
              .withIssuer(jwtDomain)
              .build())
      validate { credential ->
        if (credential.payload.audience.contains(jwtAudience)) JWTPrincipal(credential.payload)
        else null
      }
    }
  }
  routing {
    get("/session/increment") {
      val session = call.sessions.get<MySession>() ?: MySession()
      call.sessions.set(session.copy(count = session.count + 1))
      call.respondText("Counter is ${session.count}. Refresh to increment.")
    }
  }
}
