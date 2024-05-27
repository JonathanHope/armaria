Feature: Remove Folders with CLI

  The Armaria CLI can be used to remove an existing folder.
  
  @cli @remove_folder
  Scenario: Can remove folder
    Given the DB already has the following entries:
      | id            | parent_id     | is_folder | name           | url            | description | tags              |
      | {parent_1_id} | NULL          | true      | tech           |                | NULL        |                   |
      | {parent_2_id} | [parent_1_id] | true      | blogs          |                | NULL        |                   |
      | {id}          | [parent_2_id] | false     | https://jho.pe | https://jho.pe | NULL        | blog, programming |
    When I run it with the following args:
      """
      remove folder [parent_1_id]
      """
    Then the following bookmarks/folders exist:
      | id | parent_id | is_folder | name | url | description | tags |
    And the folllowing tags exist:
      | tag |

  @cli @remove_folder
  Scenario: Folder must exist
    When I run it with the following args:
      """
      remove folder test
      """
    Then the following error is returned:
      """
      Folder not found
      """
