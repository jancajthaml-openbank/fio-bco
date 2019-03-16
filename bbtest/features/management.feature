Feature: Properly behaving units

  Scenario: onboard
    Given tenant lorem is onbdoarded
    And tenant ipsum is onbdoarded
    Then systemctl contains following
    """
      fio-bco-import@lorem.service
      fio-bco-import@ipsum.service
      fio-bco-rest.service
      fio-bco.service
    """

    When stop unit "fio-bco-rest.service"
    Then unit "fio-bco-rest.service" is not running

    When start unit "fio-bco-rest.service"
    Then unit "fio-bco-rest.service" is running

    When restart unit "fio-bco-rest.service"
    Then unit "fio-bco-rest.service" is running

    When stop unit "fio-bco-import@lorem.service"
    Then unit "fio-bco-import@lorem.service" is not running

    When start unit "fio-bco-import@lorem.service"
    Then unit "fio-bco-import@lorem.service" is running

    When restart unit "fio-bco-import@ipsum.service"
    Then unit "fio-bco-import@ipsum.service" is running

  Scenario: offboard
    Given tenant lorem is offboarded
    And tenant ipsum is offboarded

    Then systemctl does not contains following
    """
      fio-bco-import@lorem.service
      fio-bco-import@ipsum.service
    """
    And systemctl contains following
    """
      fio-bco-rest.service
      fio-bco.service
    """
