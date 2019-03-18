@uninstall
Feature: Uninstall package

  Scenario: uninstall
    Given package "fio-bco" is uninstalled
    Then  systemctl does not contains following
    """
      fio-bco.service
      fio-bco.path
      fio-bco-rest.service
    """
