Feature: Messaging behaviour

  Scenario: token
    Given tenant MSG1 is onbdoarded

    When lake recieves "FioUnit/MSG1 FioRest token_1 req_id_1 NT X"
    Then lake responds with "FioRest FioUnit/MSG1 req_id_1 token_1 TN"

    When lake recieves "FioUnit/MSG1 FioRest token_2 req_id_2 NT X"
    And  lake recieves "FioUnit/MSG1 FioRest token_2 req_id_2 NT X"
    Then lake responds with "FioRest FioUnit/MSG1 req_id_2 token_2 TN"
    And  lake responds with "FioRest FioUnit/MSG1 req_id_2 token_2 EE"

    When lake recieves "FioUnit/MSG1 FioRest token_3 req_id_3 DT"
    Then lake responds with "FioRest FioUnit/MSG1 req_id_3 token_3 EE"

    When lake recieves "FioUnit/MSG1 FioRest token_4 req_id_4 NT X"
    And  lake recieves "FioUnit/MSG1 FioRest token_4 req_id_4 DT"
    Then lake responds with "FioRest FioUnit/MSG1 req_id_4 token_4 TN"
    And  lake responds with "FioRest FioUnit/MSG1 req_id_4 token_4 TD"
