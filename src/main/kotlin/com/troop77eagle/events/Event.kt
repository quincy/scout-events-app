package com.troop77eagle.events

import kotlinx.datetime.Instant
import kotlinx.serialization.Serializable

@Serializable
data class Event(
    val id: Long,
    val name: String,
    val startTime: Instant,
    val endTime: Instant,
    val summary: String,
    val description: String,
    val eventLocation: String,
    val assemblyLocation: String? = null,
    val pickupLocation: String? = null,
)

@Serializable
data class CreateEventRequest(
    val name: String,
    val startTime: String,
    val endTime: String,
    val summary: String,
    val description: String,
    val eventLocation: String,
    val assemblyLocation: String? = null,
    val pickupLocation: String? = null,
) {
  fun toEvent(id: Long): Event =
      Event(
          id = id,
          name = this.name,
          startTime = Instant.parse(this.startTime),
          endTime = Instant.parse(this.endTime),
          summary = this.summary,
          description = this.description,
          eventLocation = this.eventLocation,
          assemblyLocation = this.assemblyLocation,
          pickupLocation = this.pickupLocation)
}
