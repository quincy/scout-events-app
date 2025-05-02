Feature: add event
  Members of Troop 77 with the required permissions can create new events.

  Background:
    Given the events app is running
    And I am logged in with the "EVENT_ADMIN" role

  Scenario: The events list page has a button for adding a new event
    When I visit the events list page
    Then I should see the "Add Event" button which GETs "/events:create"

  Scenario: When I click the button, I should see the add event form
    Given no events exist
    When I visit the events list page
    And I click the "Add Event" button
    Then I should see a form POSTing to "/events:create" with the following fields:
      | Name              | Type           |
      | Name              | text           |
      | Start Time        | datetime-local |
      | End Time          | datetime-local |
      | Summary           | text           |
      | Description       | textarea       |
      | Event Location    | text           |
      | Assembly Location | text           |
      | Pickup Location   | text           |

  Scenario: When I submit the create event form, I should be redirected to the event details page for the new event
    Given no events exist
    When I visit the events list page
    And I click the "Add Event" button
    And I fill in the form with the following values:
      | Name              | Value                     |
      | Name              | campout 1                 |
      | Start Time        | tomorrow at 7pm           |
      | End Time          | tomorrow at 9pm           |
      | Summary           | camping                   |
      | Description       | fun                       |
      | Event Location    | Yellowstone National Park |
      | Assembly Location | Our meeting place         |
      | Pickup Location   | Our pickup place          |
    And I submit the form
    Then I should see the event details page with these values:
      """
      {
        "name": "campout 1",
        "start-and-end-times": "tomorrow at 7pm - tomorrow at 9pm",
        "summary": "camping",
        "description": "fun",
        "event-location": "Yellowstone National Park",
        "assembly-location": "Our meeting place",
        "pickup-location": "Our pickup place"
      }
      """
