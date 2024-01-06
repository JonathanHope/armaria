Feature: List Folder with CLI

  The Armaria CLI can be used to list folders.

  @cli @list_folders
  Scenario: Can list folders
    Given the DB already has the following entries:
      | id            | parent_id | is_folder | name           | url            | description | tags |
      | {parent_1_id} | NULL      | true      | blogs          | NULLL          | NULL        |      |
      | {parent_2_id} | NULL      | true      | tech           | NULL           | NULL        |      |
      | {id}          | NULL      | false     | https://jho.pe | https://jho.pe | NULL        |      |
    When I run it with the following args:
      """
      list folders
      """
    Then the folllowing books are returned:
      | id            | parent_id | is_folder | name  | url  | description | tags |
      | [parent_1_id] | NULL      | true      | blogs | NULL | NULL        |      |
      | [parent_2_id] | NULL      | true      | tech  | NULL | NULL        |      |

  @cli @list_folders
  Scenario: Can limit listed folders
    Given the DB already has the following entries:
      | id            | parent_id | is_folder | name           | url            | description | tags |
      | {parent_1_id} | NULL      | true      | blogs          | NULL           | NULL        |      |
      | {parent_2_id} | NULL      | true      | tech           | NULL           | NULL        |      |
      | {id}          | NULL      | false     | https://jho.pe | https://jho.pe | NULL        |      |
    When I run it with the following args:
      """
      list folders --first 1
      """
    Then the folllowing books are returned:
      | id            | parent_id | is_folder | name  | url  | description | tags |
      | [parent_1_id] | NULL      | true      | blogs | NULL | NULL        |      |

  @cli @list_folders
  Scenario: Can order folders by name ascending
    Given the DB already has the following entries:
      | id            | parent_id | is_folder | name           | url            | description | tags |
      | {parent_1_id} | NULL      | true      | blogs          | NULL           | NULL        |      |
      | {parent_2_id} | NULL      | true      | tech           | NULL           | NULL        |      |
      | {id}          | NULL      | false     | https://jho.pe | https://jho.pe | NULL        |      |
    When I run it with the following args:
      """
      list folders --order name --dir asc
      """
    Then the folllowing books are returned:
      | id            | parent_id | is_folder | name  | url  | description | tags |
      | [parent_1_id] | NULL      | true      | blogs | NULL | NULL        |      |
      | [parent_2_id] | NULL      | true      | tech  | NULL | NULL        |      |

  @cli @list_folders
  Scenario: Can order folders by name descending
    Given the DB already has the following entries:
      | id            | parent_id | is_folder | name           | url            | description | tags |
      | {parent_1_id} | NULL      | true      | blogs          | NULL           | NULL        |      |
      | {parent_2_id} | NULL      | true      | tech           | NULL           | NULL        |      |
      | {id}          | NULL      | false     | https://jho.pe | https://jho.pe | NULL        |      |
    When I run it with the following args:
      """
      list folders --order name --dir desc
      """
    Then the folllowing books are returned:
      | id            | parent_id | is_folder | name  | url  | description | tags |
      | [parent_2_id] | NULL      | true      | tech  | NULL | NULL        |      |
      | [parent_1_id] | NULL      | true      | blogs | NULL | NULL        |      |

  @cli @list_folders
  Scenario: Can order folders by manual ascending
    Given the DB already has the following entries:
      | id            | parent_id | is_folder | name           | url            | description | tags |
      | {parent_1_id} | NULL      | true      | blogs          | NULL           | NULL        |      |
      | {parent_2_id} | NULL      | true      | tech           | NULL           | NULL        |      |
      | {id}          | NULL      | false     | https://jho.pe | https://jho.pe | NULL        |      |
    When I run it with the following args:
      """
      list folders --order manual --dir asc
      """
    Then the folllowing books are returned:
      | id            | parent_id | is_folder | name  | url  | description | tags |
      | [parent_1_id] | NULL      | true      | blogs | NULL | NULL        |      |
      | [parent_2_id] | NULL      | true      | tech  | NULL | NULL        |      |

  @cli @list_folders
  Scenario: Can order folders by manual descending
    Given the DB already has the following entries:
      | id            | parent_id | is_folder | name           | url            | description | tags |
      | {parent_1_id} | NULL      | true      | blogs          | NULL           | NULL        |      |
      | {parent_2_id} | NULL      | true      | tech           | NULL           | NULL        |      |
      | {id}          | NULL      | false     | https://jho.pe | https://jho.pe | NULL        |      |
    When I run it with the following args:
      """
      list folders --order manual --dir desc
      """
    Then the folllowing books are returned:
      | id            | parent_id | is_folder | name  | url  | description | tags |
      | [parent_2_id] | NULL      | true      | tech  | NULL | NULL        |      |
      | [parent_1_id] | NULL      | true      | blogs | NULL | NULL        |      |

  @cli @list_folders
  Scenario: Can list folders after folder
    Given the DB already has the following entries:
      | id            | parent_id | is_folder | name           | url            | description | tags |
      | {parent_1_id} | NULL      | true      | blogs          | NULL           | NULL        |      |
      | {parent_2_id} | NULL      | true      | tech           | NULL           | NULL        |      |
      | {id}          | NULL      | false     | https://jho.pe | https://jho.pe | NULL        |      |
    When I run it with the following args:
      """
      list folders --order name --dir asc --after [parent_1_id]
      """
    Then the folllowing books are returned:
      | id            | parent_id | is_folder | name | url  | description | tags |
      | [parent_2_id] | NULL      | true      | tech | NULL | NULL        |      |

  @cli @list_folders
  Scenario: Can list folders in a parent folder
    Given the DB already has the following entries:
      | id            | parent_id     | is_folder | name           | url            | description | tags |
      | {parent_1_id} | NULL          | true      | blogs          | NULL           | NULL        |      |
      | {parent_2_id} | [parent_1_id] | true      | tech           | NULL           | NULL        |      |
      | {id}          | NULL          | false     | https://jho.pe | https://jho.pe | NULL        |      |
    When I run it with the following args:
      """
      list folders --folder [parent_1_id]
      """
    Then the folllowing books are returned:
      | id            | parent_id     | is_folder | name | url  | description | tags |
      | [parent_2_id] | [parent_1_id] | true      | tech | NULL | NULL        |      |

  @cli @list_folders
  Scenario: Can list top level folders
    Given the DB already has the following entries:
      | id            | parent_id     | is_folder | name           | url            | description | tags |
      | {parent_1_id} | NULL          | true      | blogs          | NULL           | NULL        |      |
      | {parent_2_id} | [parent_1_id] | true      | tech           | NULL           | NULL        |      |
      | {id}          | NULL          | false     | https://jho.pe | https://jho.pe | NULL        |      |
    When I run it with the following args:
      """
      list folders --no-folder
      """
    Then the folllowing books are returned:
      | id            | parent_id | is_folder | name  | url  | description | tags |
      | [parent_1_id] | NULL      | true      | blogs | NULL | NULL        |      |

  @cli @list_folders
  Scenario: Can search folders
    Given the DB already has the following entries:
      | id            | parent_id     | is_folder | name           | url            | description | tags |
      | {parent_1_id} | NULL          | true      | blogs          | NULL           | NULL        |      |
      | {parent_2_id} | [parent_1_id] | true      | tech           | NULL           | NULL        |      |
      | {id}          | NULL          | false     | https://jho.pe | https://jho.pe | NULL        |      |
    When I run it with the following args:
      """
      list folders --query ech
      """
    Then the folllowing books are returned:
      | id            | parent_id     | is_folder | name | url  | description | tags |
      | [parent_2_id] | [parent_1_id] | true      | tech | NULL | NULL        |      |

  @cli @list_folders
  Scenario: First must be greater than zero
    When I run it with the following args:
      """
      list folders --first 0
      """
    Then the following error is returned:
      """
      First too small
      """

  @cli @list_folders
  Scenario: Query must be at leat three chars
    When I run it with the following args:
      """
      list folders --query a
      """
    Then the following error is returned:
      """
      Query too short
      """

  @cli @list_folders
  Scenario: Cannot filter by folder and top level at same time
    When I run it with the following args:
      """
      list folders --folder [parent_id] --no-folder
      """
    Then the following error is returned:
      """
      Arguments folder and no-folder are mutually exclusive
      """
