Feature: Messaging behaviour

  Scenario: token
    Given tenant MSG is onbdoarded

    When tenant MSG receives "token req_id_1 NT token_1 token_1"
    Then tenant MSG responds with "req_id_1 token TN"
    And no other messages were received

    When tenant MSG receives "token req_id_2 NT token_2 token_2"
    And tenant MSG receives "token req_id_2 NT token_2 token_2"
    Then tenant MSG responds with "req_id_2 token TN"
    And tenant MSG responds with "req_id_2 token EE"
    And no other messages were received

    When tenant MSG receives "token req_id_3 NT token_3 token_3"
    And tenant MSG receives "token req_id_3 DT token_3"
    Then tenant MSG responds with "req_id_3 token TN"
    And tenant MSG responds with "req_id_3 token TD"
    And no other messages were received

    When tenant MSG receives "token req_id_4 DT token_4"
    Then tenant MSG responds with "req_id_4 token EE"
    And no other messages were received
