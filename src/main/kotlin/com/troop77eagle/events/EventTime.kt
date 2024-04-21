package com.troop77eagle.events

import java.time.ZoneOffset
import kotlinx.datetime.FixedOffsetTimeZone
import kotlinx.datetime.LocalTime
import kotlinx.datetime.UtcOffset
import kotlinx.serialization.Serializable

/**
 * `EventTime` represents a time in a given time zone. FIXME I want to separate time and zone info
 * and store them as separate fields in the object. In the database, we should always store UTC and
 * convert to the local time.
 */
@Serializable
data class EventTime(val time: LocalTime, val zone: FixedOffsetTimeZone) {
  constructor(
      hour: Int,
      minute: Int,
      zoneOffset: ZoneOffset,
  ) : this(LocalTime(hour, minute), FixedOffsetTimeZone(UtcOffset(zoneOffset)))

  companion object {
    /**
     * Create an `EventTime` by parsing a String representation like `HH:mm[+-]HH:mm`. Example
     * parse("13:30-08:00") -> EventTime(LocalTime(13, 30), FixedOffsetTimeZone(UtcOffset(8, 0)))
     */
    fun parse(timeStr: String): EventTime {
      val match =
          checkNotNull(Regex("""(.*)([+-].*)""").matchEntire(timeStr)) {
            "Invalid time format: $timeStr"
          }
      val time = match.groupValues[1]
      val zone = match.groupValues[2]
      return EventTime(LocalTime.parse(time), FixedOffsetTimeZone(UtcOffset(ZoneOffset.of(zone))))
    }
  }
}
