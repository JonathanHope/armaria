Feature: List Tags with CLI

  The Armaria CLI can be used to list tags.

  @cli @list_tags
  Scenario: Can list tags
    Given the DB already has the following entries:
      | id   | parent_id | is_folder | name           | url            | description | tags              |
      | {id} | NULL      | false     | https://jho.pe | https://jho.pe | NULL        | blog, programming |
    When I run it with the following args:
      """
      list tags
      """
    Then the folllowing tags are returned:
      | tag         |
      | blog        |
      | programming |

  @cli @list_tags
  Scenario: Can order tags by tag ascending
    Given the DB already has the following entries:
      | id   | parent_id | is_folder | name           | url            | description | tags              |
      | {id} | NULL      | false     | https://jho.pe | https://jho.pe | NULL        | blog, programming |
    When I run it with the following args:
      """
      list tags --dir asc
      """
    Then the folllowing tags are returned:
      | tag         |
      | blog        |
      | programming |

  @cli @list_tags
  Scenario: Can order tags by tag descending
    Given the DB already has the following entries:
      | id   | parent_id | is_folder | name           | url            | description | tags              |
      | {id} | NULL      | false     | https://jho.pe | https://jho.pe | NULL        | blog, programming |
    When I run it with the following args:
      """
      list tags --dir desc
      """
    Then the folllowing tags are returned:
      | tag         |
      | programming |
      | blog        |

  @cli @list_tags
  Scenario: Can limit listed tags
    Given the DB already has the following entries:
      | id   | parent_id | is_folder | name           | url            | description | tags              |
      | {id} | NULL      | false     | https://jho.pe | https://jho.pe | NULL        | blog, programming |
    When I run it with the following args:
      """
      list tags --first 1
      """
    Then the folllowing tags are returned:
      | tag  |
      | blog |

  @cli @list_tags
  Scenario: Can list tags after a tag
    Given the DB already has the following entries:
      | id   | parent_id | is_folder | name           | url            | description | tags              |
      | {id} | NULL      | false     | https://jho.pe | https://jho.pe | NULL        | blog, programming |
    When I run it with the following args:
      """
      list tags --after blog
      """
    Then the folllowing tags are returned:
      | tag         |
      | programming |

  @cli @list_tags
  Scenario: Can search tags
    Given the DB already has the following entries:
      | id   | parent_id | is_folder | name           | url            | description | tags              |
      | {id} | NULL      | false     | https://jho.pe | https://jho.pe | NULL        | blog, programming |
    When I run it with the following args:
      """
      list tags --query gram
      """
    Then the folllowing tags are returned:
      | tag         |
      | programming |

  @cli @list_tags
  Scenario: First must be greater than zero
    When I run it with the following args:
      """
      list tags --first 0
      """
    Then the following error is returned:
      """
      First too small
      """

  @cli @list_tags
  Scenario: Query must be at leat three chars
    When I run it with the following args:
      """
      list tags --query a
      """
    Then the following error is returned:
      """
      Query too short
      """
