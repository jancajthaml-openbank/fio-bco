Feature: Properly behaving units

  Scenario: onboard
    Given tenant lorem is onboarded
    And   tenant ipsum is onboarded
    Then  systemctl contains following active units
      | name                 | type    |
      | fio-bco              | path    |
      | fio-bco              | service |
      | fio-bco-rest         | service |
      | fio-bco-import@lorem | service |
      | fio-bco-import@ipsum | service |
    And unit "fio-bco-import@lorem.service" is running
    And unit "fio-bco-import@ipsum.service" is running

    When stop unit "fio-bco-import@lorem.service"
    Then unit "fio-bco-import@lorem.service" is not running
    And  unit "fio-bco-import@ipsum.service" is running

    When start unit "fio-bco-import@lorem.service"
    Then unit "fio-bco-import@lorem.service" is running

    When restart unit "fio-bco-import@lorem.service"
    Then unit "fio-bco-import@lorem.service" is running

  Scenario: offboard
    Given tenant lorem is offboarded
    And   tenant ipsum is offboarded
    Then  systemctl does not contain following active units
      | name                 | type    |
      | fio-bco-import@lorem | service |
      | fio-bco-import@ipsum | service |
    And systemctl contains following active units
      | name                 | type    |
      | fio-bco              | path    |
      | fio-bco              | service |
      | fio-bco-rest         | service |
