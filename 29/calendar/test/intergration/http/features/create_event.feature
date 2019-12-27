Feature: Create event
  As API client of calendar service
  I want to create new event

  Scenario: Create event, 200 OK
    Given Clean DB
    When I send "POST" request to "http://localhost:8888/create_event" with "application/x-www-form-urlencoded" params:
    """
    name=Add test&start=2019-12-21 14:00&end=2019-12-21 15:00&beforeMinutes=10
    """
    Then The response code should be 200
    And The response contentType should be "application/json"
    And The response json should has field "result" with value match "^created (\d+)$"
    And Extracted number is event id
    And The record should match:
      | name      | start_time        | end_time          | before_minutes  | notified_time |
      | Add test  | 2019-12-21 14:00  | 2019-12-21 15:00  | 10              | nil           |

  Scenario: Create event, 400 invalid start date
    Given Clean DB
    When I send "POST" request to "http://localhost:8888/create_event" with "application/x-www-form-urlencoded" params:
    """
    name=Add test&start=sdfasdf&end=2019-12-21 15:00&beforeMinutes=10
    """
    Then The response code should be 400
    And The response contentType should be "application/json"
    And The response json should has field "error" with value match "^invalid format of datetime"

  Scenario: Create event, 400 invalid end date
    Given Clean DB
    When I send "POST" request to "http://localhost:8888/create_event" with "application/x-www-form-urlencoded" params:
    """
    name=Add test&start=2019-12-21 14:00&end=dfwedfe&beforeMinutes=10
    """
    Then The response code should be 400
    And The response contentType should be "application/json"
    And The response json should has field "error" with value match "^invalid format of datetime"
