Feature: This will test with some table for a step

    Scenario: First Scenario
        Given I have something
        When I do something
        Then something happened
            | name | value |
            | foo  | bar   |
            | hello| world |

    Scenario: Second Scenario
        Given I have something else
        When I do something different
        Then something else happened

    Scenario Outline: Third Scenario
        Given I have something
        When I do something
        Then something happened
            | name | value |
            | <key>  | <value>   |
            | hello| world |
        And something else happened

    Examples: Some exmaples
        | key | value |
        | foo | bar |
        | lucky | luke |

    Scenario: Fourth scenario
        Given I have something
        When I do something
        Then something happened
            And something different happened
                | key | value |
                | foo | bar |
                | hello | there |

    Scenario: Fifth scenario
        Given I have something
        When I do something
        Then something happened
