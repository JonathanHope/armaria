Feature: List Books with CLI

  The Armaria CLI can be used to list books.

  @cli @list_books
  Scenario: Can list bookmarks
    Given the DB already has the following entries:
      | id          | parent_id | is_folder | name                | url                 | description | tags |
      | {parent_id} | NULL      | true      | blogs               | NULL                | NULL        |      |
      | {id_1}      | NULL      | false     | https://jho.pe      | https://jho.pe      | NULL        |      |
      | {id_2}      | NULL      | false     | https://armaria.net | https://armaria.net | NULL        |      |
    When I run it with the following args:
      """
      list books
      """
    Then the folllowing books are returned:
      | id     | parent_id | is_folder | name                | url                 | description | tags |
      | [id_1] | NULL      | false     | https://jho.pe      | https://jho.pe      | NULL        |      |
      | [id_2] | NULL      | false     | https://armaria.net | https://armaria.net | NULL        |      |

  @cli @list_books
  Scenario: Can limit listed bookmarks
    Given the DB already has the following entries:
      | id          | parent_id | is_folder | name                | url                 | description | tags |
      | {parent_id} | NULL      | true      | blogs               | NULL                | NULL        |      |
      | {id_1}      | NULL      | false     | https://jho.pe      | https://jho.pe      | NULL        |      |
      | {id_2}      | NULL      | false     | https://armaria.net | https://armaria.net | NULL        |      |
    When I run it with the following args:
      """
      list books --first 1
      """
    Then the folllowing books are returned:
      | id     | parent_id | is_folder | name           | url            | description | tags |
      | [id_1] | NULL      | false     | https://jho.pe | https://jho.pe | NULL        |      |

  @cli @list_books
  Scenario: Can order bookmarks by name ascending
    Given the DB already has the following entries:
      | id          | parent_id | is_folder | name                | url                 | description | tags |
      | {parent_id} | NULL      | true      | blogs               | NULL                | NULL        |      |
      | {id_1}      | NULL      | false     | https://jho.pe      | https://jho.pe      | NULL        |      |
      | {id_2}      | NULL      | false     | https://armaria.net | https://armaria.net | NULL        |      |
    When I run it with the following args:
      """
      list books --order name --dir asc
      """
    Then the folllowing books are returned:
      | id     | parent_id | is_folder | name                | url                 | description | tags |
      | [id_2] | NULL      | false     | https://armaria.net | https://armaria.net | NULL        |      |
      | [id_1] | NULL      | false     | https://jho.pe      | https://jho.pe      | NULL        |      |

  @cli @list_books
  Scenario: Can order bookmarks by name descending
    Given the DB already has the following entries:
      | id          | parent_id | is_folder | name                | url                 | description | tags |
      | {parent_id} | NULL      | true      | blogs               | NULL                | NULL        |      |
      | {id_1}      | NULL      | false     | https://jho.pe      | https://jho.pe      | NULL        |      |
      | {id_2}      | NULL      | false     | https://armaria.net | https://armaria.net | NULL        |      |
    When I run it with the following args:
      """
      list books --order name --dir desc
      """
    Then the folllowing books are returned:
      | id     | parent_id | is_folder | name                | url                 | description | tags |
      | [id_1] | NULL      | false     | https://jho.pe      | https://jho.pe      | NULL        |      |
      | [id_2] | NULL      | false     | https://armaria.net | https://armaria.net | NULL        |      |

  @cli @list_books
  Scenario: Can list bookmarks after bookmark
    Given the DB already has the following entries:
      | id          | parent_id | is_folder | name                | url                 | description | tags |
      | {parent_id} | NULL      | true      | blogs               | NULL                | NULL        |      |
      | {id_1}      | NULL      | false     | https://jho.pe      | https://jho.pe      | NULL        |      |
      | {id_2}      | NULL      | false     | https://armaria.net | https://armaria.net | NULL        |      |
    When I run it with the following args:
      """
      list books --order name --dir asc --after [id_2]
      """
    Then the folllowing books are returned:
      | id     | parent_id | is_folder | name           | url            | description | tags |
      | [id_1] | NULL      | false     | https://jho.pe | https://jho.pe | NULL        |      |

  @cli @list_books
  Scenario: Can list bookmarks in a parent folder
    Given the DB already has the following entries:
      | id          | parent_id   | is_folder | name                | url                 | description | tags |
      | {parent_id} | NULL        | true      | blogs               | NULL                | NULL        |      |
      | {id_1}      | [parent_id] | false     | https://jho.pe      | https://jho.pe      | NULL        |      |
      | {id_2}      | NULL        | false     | https://armaria.net | https://armaria.net | NULL        |      |
    When I run it with the following args:
      """
      list books --folder [parent_id]
      """
    Then the folllowing books are returned:
      | id     | parent_id   | is_folder | name           | url            | description | tags |
      | [id_1] | [parent_id] | false     | https://jho.pe | https://jho.pe | NULL        |      |

  @cli @list_books
  Scenario: Can search bookmarks
    Given the DB already has the following entries:
      | id          | parent_id   | is_folder | name                | url                 | description | tags |
      | {parent_id} | NULL        | true      | blogs               | NULL                | NULL        |      |
      | {id_1}      | [parent_id] | false     | https://jho.pe      | https://jho.pe      | NULL        |      |
      | {id_2}      | NULL        | false     | https://armaria.net | https://armaria.net | NULL        |      |
    When I run it with the following args:
      """
      list books --query jho
      """
    Then the folllowing books are returned:
      | id     | parent_id   | is_folder | name           | url            | description | tags |
      | [id_1] | [parent_id] | false     | https://jho.pe | https://jho.pe | NULL        |      |

  @cli @list_books
  Scenario: Can filter bookmarks by tag
    Given the DB already has the following entries:
      | id          | parent_id | is_folder | name                | url                 | description | tags |
      | {parent_id} | NULL      | true      | blogs               | NULL                | NULL        |      |
      | {id_1}      | NULL      | false     | https://jho.pe      | https://jho.pe      | NULL        | blog |
      | {id_2}      | NULL      | false     | https://armaria.net | https://armaria.net | NULL        |      |
    When I run it with the following args:
      """
      list books --tag blog
      """
    Then the folllowing books are returned:
      | id     | parent_id | is_folder | name           | url            | description | tags |
      | [id_1] | NULL      | false     | https://jho.pe | https://jho.pe | NULL        | blog |

  @cli @list_books
  Scenario: First must be greater than zero
    When I run it with the following args:
      """
      list books --first 0
      """
    Then the following error is returned:
      """
      First too small
      """

  @cli @list_books
  Scenario: Query must be at leat three chars
    When I run it with the following args:
      """
      list books --query a
      """
    Then the following error is returned:
      """
      Query too short
      """
