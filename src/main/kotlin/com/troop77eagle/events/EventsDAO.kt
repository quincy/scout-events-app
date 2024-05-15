package com.troop77eagle.events

import java.sql.ResultSet
import kotlinx.datetime.Instant
import kotlinx.datetime.UtcOffset
import kotlinx.datetime.format.DateTimeComponents
import kotlinx.datetime.format.char
import org.jdbi.v3.core.mapper.ColumnMapper
import org.jdbi.v3.core.statement.StatementContext
import org.jdbi.v3.sqlobject.statement.GetGeneratedKeys
import org.jdbi.v3.sqlobject.statement.SqlQuery
import org.jdbi.v3.sqlobject.statement.SqlUpdate

interface EventsDAO {

  @SqlQuery("""select * from scouting.public.events where id = :id""") fun getById(id: Long): Event?

  @SqlUpdate(
      """insert into scouting.public.events
           (name,
            start_time,
            end_time,
            summary,
            description,
            event_location,
            assembly_location,
            pickup_location)
         values (:request.name,
                 :request.startTime,
                 :request.endTime,
                 :request.summary,
                 :request.description,
                 :request.eventLocation,
                 :request.assemblyLocation,
                 :request.pickupLocation)""")
  @GetGeneratedKeys("id")
  fun create(request: CreateEventRequest): Long
}

class InstantMapper : ColumnMapper<Instant> {
  // "yyyy-MM-dd HH:mm:ssZ"
  private val format =
      DateTimeComponents.Format {
        year()
        char('-')
        monthNumber()
        char('-')
        dayOfMonth()
        char(' ')
        hour()
        char(':')
        minute()
        char(':')
        second()
        offset(UtcOffset.Formats.ISO_BASIC)
      }

  override fun map(r: ResultSet, columnNumber: Int, ctx: StatementContext): Instant {
    return Instant.parse(r.getString(columnNumber), format)
  }
}
