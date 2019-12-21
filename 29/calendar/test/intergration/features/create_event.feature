Feature: Create event
  As API client of calendar service
  I want to create new event

  Scenario: Create event
    When I send "POST" request to "http://localhost:8888/create_event" with "application/x-www-form-urlencoded" params:
    """
    name=Add test&start=2019-12-21 14:00&end=2019-12-21 15:00&beforeMinutes=10
    """
    Then The response code should be 200
    And The response contentType should be "application/json"
    And The response json should has field "result" with value match "^created \d+$"
