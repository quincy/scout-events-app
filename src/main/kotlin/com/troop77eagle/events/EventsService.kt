package com.troop77eagle.events

import io.github.oshai.kotlinlogging.KotlinLogging
import org.jdbi.v3.core.JdbiException

interface EventsService {
  fun create(request: CreateEventRequest): Event

  fun getById(id: Long): Event?

  fun getAll(): List<Event>

  companion object {
    fun create(eventsDAO: EventsDAO): EventsService {
      return EventsServiceImpl(eventsDAO)
    }
  }
}

private val log = KotlinLogging.logger {}

private class EventsServiceImpl(private val eventsDAO: EventsDAO) : EventsService {
  override fun create(request: CreateEventRequest): Event {
    return eventsDAO.create(request).let { request.toEvent(it) }
  }

  override fun getById(id: Long): Event? {
    return eventsDAO.getById(id)
  }

  override fun getAll(): List<Event> {
    return try {
      eventsDAO.getAll()
    } catch (e: JdbiException) {
      log.error(e) { "SQL query failed" }
      throw e
    }
  }
}
