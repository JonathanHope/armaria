Feature: Get Bookmark or Folder with CLI

  The Armaria CLI can be used to add tags to an existing bookmark.

  @cli @get_book
  Scenario: Can get a bookmark
    Given the DB already has the following entries:
      | id   | parent_id | is_folder | name           | url            | description | tags |
      | {id} | NULL      | false     | https://jho.pe | https://jho.pe | NULL        |      |
    When I run it with the following args:
      """
      get all [id]
      """
    Then the folllowing books are returned:
      | id   | parent_id | is_folder | name           | url            | description | tags |
      | [id] | NULL      | false     | https://jho.pe | https://jho.pe | NULL        |      |

  @cli @get_book
  Scenario: Can get a bookmark
    Given the DB already has the following entries:
      | id   | parent_id | is_folder | name  | url  | description | tags |
      | {id} | NULL      | true      | blogs | NULL | NULL        |      |
    When I run it with the following args:
      """
      get all [id]
      """
    Then the folllowing books are returned:
      | id   | parent_id | is_folder | name  | url  | description | tags |
      | [id] | NULL      | true      | blogs | NULL | NULL        |      |

  @cli @get_book
  Scenario: Bookmark or folder must exist
    When I run it with the following args:
      """
      get all test
      """
    Then the following error is returned:
      """
      Bookmark or folder not found
      """
