Feature: Query with CLI

  The Armaria CLI can be used to query bookmarks and folders.

  @cli @query
  Scenario: Can query bookmarks/folders
    Given the DB already has the following entries:
      | id            | parent_id | is_folder | name           | url            | description | tags |
      | {parent_1_id} | NULL      | true      | blogs          | NULL           | NULL        |      |
      | {id}          | NULL      | false     | https://jho.pe | https://jho.pe | NULL        |      |
    When I run it with the following args:
      """
      query blog
      """
    Then the folllowing books are returned:
      | id            | parent_id | is_folder | name           | url            | description | tags |
      | [parent_1_id] | NULL      | true      | blogs          | NULL           | NULL        |      |

  @cli @query
  Scenario: Can limit query results
    Given the DB already has the following entries:
      | id            | parent_id | is_folder | name           | url            | description | tags |
      | {parent_1_id} | NULL      | true      | blogs          | NULL           | NULL        |      |
      | {parent_2_id} | NULL      | true      | blogs          | NULL           | NULL        |      |
    When I run it with the following args:
      """
      query blog --first 1
      """
    Then the folllowing books are returned:
      | id            | parent_id | is_folder | name           | url            | description | tags |
      | [parent_1_id] | NULL      | true      | blogs          | NULL           | NULL        |      |
