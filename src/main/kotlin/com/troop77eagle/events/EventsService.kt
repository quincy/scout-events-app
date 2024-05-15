package com.troop77eagle.events

interface EventsService {
  fun create(request: CreateEventRequest): Event

  fun getById(id: Long): Event?

  companion object {
    fun create(eventsDAO: EventsDAO): EventsService {
      return EventsServiceImpl(eventsDAO)
    }
  }
}

private class EventsServiceImpl(private val eventsDAO: EventsDAO) : EventsService {
  override fun create(request: CreateEventRequest): Event {
    return eventsDAO.create(request).let { request.toEvent(it) }
  }

  override fun getById(id: Long): Event? {
    return eventsDAO.getById(id)
  }
}
