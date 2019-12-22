Feature: Update event
  As API client of calendar service
  I want to update existing event

  Scenario: Update event
    Given Existing records:
      | name       | start_time        | end_time          | before_minutes  | notified_time |
      | test1      | 2019-12-21 14:00  | 2019-12-21 15:00  | 10              | nil           |
      | test2      | 2019-12-22 15:00  | 2019-12-22 17:00  | nil             | nil           |

    When I send POST request to "http://localhost:8888/update_event" with "application/x-www-form-urlencoded" params:
    # about id=1 before,
    # actually "1" here is index of record in Given section table data; In test "1" will be replaced by real event ID
    # because we will know about real IDs only after inserting records in DB
    """
    id=1&name=updated&start=2019-12-22 15:15&end=2019-12-22 18:00&beforeMinutes=5
    """
    Then The response code should be 200
    And The response contentType should be "application/json"
    And The response json should has field "result" with value match "updated"
    And The records should match:
      | no | name       | start_time        | end_time          | before_minutes  | notified_time |
      | 1  | test1      | 2019-12-21 14:00  | 2019-12-21 15:00  | 10              | nil           |
      | 2  | updated    | 2019-12-22 15:15  | 2019-12-22 18:00  | 5               | nil           |