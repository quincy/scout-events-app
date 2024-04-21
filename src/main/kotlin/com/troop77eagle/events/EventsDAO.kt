package com.troop77eagle.events

import org.jdbi.v3.sqlobject.statement.GetGeneratedKeys
import org.jdbi.v3.sqlobject.statement.SqlUpdate

interface EventsDAO {
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
