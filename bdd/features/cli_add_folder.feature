Feature: Add Folders with CLI

  The Armaria CLI can be used to add a new folder.

  @cli @add_folder
  Scenario: Can add a folder
    When I run it with the following args:
      """
      add folder blogs
      """
    Then the following bookmarks/folders exist:
      | id   | parent_id | is_folder | name  | url  | description | tags |
      | {id} | NULL      | true      | blogs | NULL | NULL        |      |

  @cli @add_folder
  Scenario: Can add a folder to a folder
    Given the DB already has the following entries:
      | id          | parent_id | is_folder | name | url  | description | tags |
      | {parent_id} | NULL      | true      | tech | NULL | NULL        |      |
    When I run it with the following args:
      """
      add folder blogs --folder [parent_id]
      """
    Then the following bookmarks/folders exist:
      | id          | parent_id   | is_folder | name  | url  | description | tags |
      | [parent_id] | NULL        | true      | tech  | NULL | NULL        |      |
      | {id}        | [parent_id] | true      | blogs | NULL | NULL        |      |

  @cli @add_folder
  Scenario: Folder must exist
    When I run it with the following args:
      """
      add folder blogs --folder test
      """
    Then the following error is returned:
      """
      Folder not found
      """

  @cli @add_folder
  Scenario: Name must be at least 1 char
    When I run it with the following args:
      """
      add folder ""
      """
    Then the following error is returned:
      """
      Name too short
      """

  @cli @add_folder
  Scenario: Name must be at most 2048 chars
    When I run it with the following args:
      """
      add folder %repeat:x:2049%
      """
    Then the following error is returned:
      """
      Name too long
      """
