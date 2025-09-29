Feature: Feature with background

    Background:
      Given Same setup as every year
    
    Scenario: When I do something things should change
      When I change something
      Then things have changed

    Scenario: When I do nothing things should not change
      When I change nothing
      Then thing have not changed