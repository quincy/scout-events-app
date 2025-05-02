Feature: events list
  Members of Troop 77 need to be able to view the list of upcoming events.

  Background:
    Given the events app is running
    And I am logged in

  Scenario: When events list is empty, show no upcoming events message
    Given no events exist
    When I visit the events list page
    Then I should see "No upcoming events"

  Scenario: When one event exists, show the list of events
    Given no events exist
    And the following events have been created:
      """
      [
        {
          "id": "1",
          "name": "campout 1",
          "start_time": "tomorrow at 7pm",
          "end_time": "tomorrow at 9pm",
          "summary": "camping",
          "description": "fun",
          "event_location": "Yellowstone National Park",
          "assembly_location": "Our meeting place",
          "pickup_location": "Our pickup place"
        }
      ]
      """
    When I visit the events list page
    Then I should see "campout 1"
    And I should see "camping"
    And I should see "Yellowstone National Park"

  Scenario: When I click on an event, I should see the event details
    Given no events exist
    And the following events have been created:
      """
      [
        {
          "id": "1",
          "name": "campout 1",
          "start_time": "tomorrow at 7pm",
          "end_time": "tomorrow at 9pm",
          "summary": "camping",
          "description": "fun",
          "event_location": "Yellowstone National Park",
          "assembly_location": "Our meeting place",
          "pickup_location": "Our pickup place"
        }
      ]
      """
    When I visit the events list page
    And I click the "campout 1" event row
    Then I should see "camping"
    And I should see "Yellowstone National Park"
    And I should see "fun"
    And I should see "Our meeting place"
    And I should see "Our pickup place"
