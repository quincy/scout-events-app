package com.troop77eagle.events

import com.github.mustachejava.DefaultMustacheFactory
import com.github.mustachejava.MustacheFactory
import java.io.StringWriter

interface Templator {
  companion object {
    fun create(): Templator {
      return EventsTemplator(DefaultMustacheFactory())
    }
  }

  fun toHtmlTable(events: Collection<Event>): String
}

class EventsTemplator(private val mustacheFactory: MustacheFactory) : Templator {
  private val eventsListHtmlTable =
      mustacheFactory.compile(
          this::class.java.getResourceAsStream("/templates/events-table.mustache")!!.reader(),
          "events-table-html")

  override fun toHtmlTable(events: Collection<Event>): String =
      StringWriter()
          .also { writer -> eventsListHtmlTable.execute(writer, events).flush() }
          .toString()
}
