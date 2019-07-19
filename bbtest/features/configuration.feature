Feature: Service can be configured

  Scenario: configure log level to DEBUG
    Given tenant CONFIGURATION_DEBUG is onbdoarded
    And fio-bco is reconfigured with
    """
      LOG_LEVEL=DEBUG
    """
    Then journalctl of "fio-bco-import@CONFIGURATION_DEBUG.service" contains following
    """
      Log level set to DEBUG
    """

  Scenario: configure log level to ERROR
    Given tenant CONFIGURATION_ERROR is onbdoarded
    And fio-bco is reconfigured with
    """
      LOG_LEVEL=ERROR
    """
    Then journalctl of "fio-bco-import@CONFIGURATION_ERROR.service" contains following
    """
      Log level set to ERROR
    """

  Scenario: configure log level to INFO
    Given tenant CONFIGURATION_INFO is onbdoarded
    And fio-bco is reconfigured with
    """
      LOG_LEVEL=INFO
    """
    Then journalctl of "fio-bco-import@CONFIGURATION_INFO.service" contains following
    """
      Log level set to INFO
    """
