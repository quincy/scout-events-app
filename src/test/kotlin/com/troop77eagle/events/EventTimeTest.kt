package com.troop77eagle.events

import io.kotest.matchers.shouldBe
import java.time.ZoneOffset
import java.util.stream.Stream
import org.junit.jupiter.params.ParameterizedTest
import org.junit.jupiter.params.provider.Arguments
import org.junit.jupiter.params.provider.MethodSource

class EventTimeTest {
  companion object {
    @JvmStatic
    fun parseInputs(): Stream<Arguments> {
      return Stream.of(
          Arguments.of("01:00+01:00", EventTime(1, 0, ZoneOffset.ofHours(1))),
          Arguments.of("01:00-01:00", EventTime(1, 0, ZoneOffset.ofHours(-1))),
          Arguments.of("13:00+08:00", EventTime(13, 0, ZoneOffset.ofHours(8))),
          Arguments.of("13:30-08:00", EventTime(13, 30, ZoneOffset.ofHours(-8))),
          Arguments.of("19:15-07:00", EventTime(19, 15, ZoneOffset.ofHours(-7))),
      )
    }
  }

  @ParameterizedTest
  @MethodSource(value = ["parseInputs"])
  fun `test parse`(input: String, expected: EventTime) {
    EventTime.parse(input) shouldBe expected
  }
}
