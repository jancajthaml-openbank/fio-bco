Feature: Properly behaving units

  Scenario: onboard
    Given tenant lorem is onbdoarded
    And tenant ipsum is onbdoarded
    Then systemctl contains following
    """
      fio-bco.service
      fio-bco@lorem.service
      fio-bco@ipsum.service
    """

    When stop unit "fio-bco.service"
    Then unit "fio-bco.service" is not running

    When start unit "fio-bco.service"
    Then unit "fio-bco.service" is running

    When restart unit "fio-bco.service"
    Then unit "fio-bco.service" is running

    When stop unit "fio-bco@lorem.service"
    Then unit "fio-bco@lorem.service" is not running

    When start unit "fio-bco@lorem.service"
    Then unit "fio-bco@lorem.service" is running

    When restart unit "fio-bco@ipsum.service"
    Then unit "fio-bco@ipsum.service" is running

  Scenario: offboard
    Given tenant lorem is offboarded
    And tenant ipsum is offboarded

    Then systemctl does not contains following
    """
      fio-bco@lorem.service
      fio-bco@ipsum.service
    """
    And systemctl contains following
    """
      fio-bco.service
    """
