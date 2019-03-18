@install
Feature: Install package

  Scenario: install
    Given package "fio-bco.deb" is installed
    Then  systemctl contains following
    """
      fio-bco.service
      fio-bco.path
      fio-bco-rest.service
    """
