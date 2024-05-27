Feature: Remove Books with CLI

  The Armaria CLI can be used to remove an existing book.
  
  @cli @remove_book
  Scenario: Can remove bookmark
    Given the DB already has the following entries:
      | id   | parent_id | is_folder | name           | url            | description | tags              |
      | {id} | NULL      | false     | https://jho.pe | https://jho.pe | NULL        | blog, programming |
    When I run it with the following args:
      """
      remove book [id]
      """
    Then the following bookmarks/folders exist:
      | id | parent_id | is_folder | name | url | description |
    And the folllowing tags exist:
      | tag |

  @cli @remove_book
  Scenario: Bookmark must exist
    When I run it with the following args:
      """
      remove book test
      """
    Then the following error is returned:
      """
      Bookmark not found
      """
