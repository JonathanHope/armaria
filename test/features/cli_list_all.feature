Feature: List All with CLI

  The Armaria CLI can be used to list bookmarks and folders.

  @cli @list_all
  Scenario: Can list bookmarks/folders
    Given the DB already has the following entries:
      | id            | parent_id | is_folder | name           | url            | description | tags |
      | {parent_1_id} | NULL      | true      | blogs          | NULL           | NULL        |      |
      | {id}          | NULL      | false     | https://jho.pe | https://jho.pe | NULL        |      |
    When I run it with the following args:
      """
      list all
      """
    Then the folllowing books are returned:
      | id            | parent_id | is_folder | name           | url            | description | tags |
      | [parent_1_id] | NULL      | true      | blogs          | NULL           | NULL        |      |
      | [id]          | NULL      | false     | https://jho.pe | https://jho.pe | NULL        |      |

  @cli @list_all
  Scenario: Can limit listed bookmarks/folders
  Given the DB already has the following entries:
    | id            | parent_id | is_folder | name           | url            | description | tags |
    | {parent_1_id} | NULL      | true      | blogs          | NULL           | NULL        |      |
    | {id}          | NULL      | false     | https://jho.pe | https://jho.pe | NULL        |      |
  When I run it with the following args:
    """
    list all --first 1
    """
  Then the folllowing books are returned:
    | id            | parent_id | is_folder | name  | url  | description | tags |
    | [parent_1_id] | NULL      | true      | blogs | NULL | NULL        |      |

  @cli @list_all
  Scenario: Can order bookmarks/folders by name ascending
    Given the DB already has the following entries:
      | id            | parent_id | is_folder | name           | url            | description | tags |
      | {parent_1_id} | NULL      | true      | blogs          | NULL           | NULL        |      |
      | {id}          | NULL      | false     | https://jho.pe | https://jho.pe | NULL        |      |
    When I run it with the following args:
      """
      list all --order name --dir asc
      """
    Then the folllowing books are returned:
      | id            | parent_id | is_folder | name           | url            | description | tags |
      | [parent_1_id] | NULL      | true      | blogs          | NULL           | NULL        |      |
      | [id]          | NULL      | false     | https://jho.pe | https://jho.pe | NULL        |      |

  @cli @list_all
  Scenario: Can order bookmarks/folders by name descending
    Given the DB already has the following entries:
      | id            | parent_id | is_folder | name           | url            | description | tags |
      | {parent_1_id} | NULL      | true      | blogs          | NULL           | NULL        |      |
      | {id}          | NULL      | false     | https://jho.pe | https://jho.pe | NULL        |      |
    When I run it with the following args:
      """
      list all --order name --dir desc
      """
    Then the folllowing books are returned:
      | id            | parent_id | is_folder | name           | url            | description | tags |
      | [id]          | NULL      | false     | https://jho.pe | https://jho.pe | NULL        |      |
      | [parent_1_id] | NULL      | true      | blogs          | NULL           | NULL        |      |

  @cli @list_all
  Scenario: Can list bookmarks/folders after bookmark/folder
    Given the DB already has the following entries:
      | id            | parent_id | is_folder | name           | url            | description | tags |
      | {parent_1_id} | NULL      | true      | blogs          | NULL           | NULL        |      |
      | {id}          | NULL      | false     | https://jho.pe | https://jho.pe | NULL        |      |
    When I run it with the following args:
      """
      list all --order name --dir asc --after [parent_1_id]
      """
    Then the folllowing books are returned:
      | id   | parent_id | is_folder | name           | url            | description | tags |
      | [id] | NULL      | false     | https://jho.pe | https://jho.pe | NULL        |      |

  @cli @list_all
  Scenario: Can list bookmarks/folders in a parent folder
    Given the DB already has the following entries:
      | id            | parent_id     | is_folder | name           | url            | description | tags |
      | {parent_1_id} | NULL          | true      | blogs          | NULL           | NULL        |      |
      | {id}          | [parent_1_id] | false     | https://jho.pe | https://jho.pe | NULL        |      |
  When I run it with the following args:
    """
    list all --folder [parent_1_id]
    """
  Then the folllowing books are returned:
    | id   | parent_id     | is_folder | name           | url            | description | tags |
    | [id] | [parent_1_id] | false     | https://jho.pe | https://jho.pe | NULL        |      |

  @cli @list_all
  Scenario: Can list top level bookmarks/folders
    Given the DB already has the following entries:
      | id            | parent_id     | is_folder | name           | url            | description | tags |
      | {parent_1_id} | NULL          | true      | blogs          | NULL           | NULL        |      |
      | {id}          | [parent_1_id] | false     | https://jho.pe | https://jho.pe | NULL        |      |
  When I run it with the following args:
    """
    list all --no-folder
    """
  Then the folllowing books are returned:
    | id            | parent_id | is_folder | name  | url  | description | tags |
    | [parent_1_id] | NULL      | true      | blogs | NULL | NULL        |      |

  @cli @list_all
  Scenario: Can search bookmarks/folders
    Given the DB already has the following entries:
      | id            | parent_id     | is_folder | name           | url            | description | tags |
      | {parent_1_id} | NULL          | true      | blogs          |                | NULL        |      |
      | {id}          | [parent_1_id] | false     | https://jho.pe | https://jho.pe | NULL        |      |
    When I run it with the following args:
      """
      list all --query jho
      """
    Then the folllowing books are returned:
      | id   | parent_id     | is_folder | name           | url            | description | tags |
      | [id] | [parent_1_id] | false     | https://jho.pe | https://jho.pe | NULL        |      |

  @cli @list_all
  Scenario: Can filter bookmarks/folders by tag
    Given the DB already has the following entries:
      | id            | parent_id     | is_folder | name           | url            | description | tags |
      | {parent_1_id} | NULL          | true      | blogs          | NULL           | NULL        |      |
      | {id}          | [parent_1_id] | false     | https://jho.pe | https://jho.pe | NULL        | blog |
    When I run it with the following args:
      """
      list all --tag blog
      """
    Then the folllowing books are returned:
      | id   | parent_id     | is_folder | name           | url            | description | tags |
      | [id] | [parent_1_id] | false     | https://jho.pe | https://jho.pe | NULL        | blog |

  @cli @list_all
  Scenario: First must be greater than zero
    When I run it with the following args:
      """
      list all --first 0
      """
    Then the following error is returned:
      """
      First too small
      """

  @cli @list_all
  Scenario: Query must be at leat three chars
    When I run it with the following args:
      """
      list all --query a
      """
    Then the following error is returned:
      """
    Query too short
      """

  @cli @list_all
  Scenario: Cannot filter by folder and top level at same time
    When I run it with the following args:
      """
      list all --folder [parent_id] --no-folder
      """
    Then the following error is returned:
      """
      Arguments folder and no-folder are mutually exclusive
      """
