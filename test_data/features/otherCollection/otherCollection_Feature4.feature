Feature: Some feature with Examples

    Scenario Outline: Some scenario with different examples
        Given I have something
        When I do <action>
        Then I get a good result
    
    Examples: Example 1
        | action |
        | greet |
        | compliment |

    Examples: Example 2
        | action |
        | not talk bad |
        | not ignore someone |