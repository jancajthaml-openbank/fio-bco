Feature: Messaging behaviour

  Scenario: token
    Given tenant MSG1 is onboarded

    When lake recieves "FioImport/MSG1 FioRest token_1 req_id_1 NT X"
    Then lake responds with "FioRest FioImport/MSG1 req_id_1 token_1 TN"

    When lake recieves "FioImport/MSG1 FioRest token_2 req_id_2 DT"
    Then lake responds with "FioRest FioImport/MSG1 req_id_2 token_2 EM"

    When lake recieves "FioImport/MSG1 FioRest token_3 req_id_3 NT X"
    Then lake responds with "FioRest FioImport/MSG1 req_id_3 token_3 TN"
    When lake recieves "FioImport/MSG1 FioRest token_3 req_id_3 DT"
    Then lake responds with "FioRest FioImport/MSG1 req_id_3 token_3 TD"

