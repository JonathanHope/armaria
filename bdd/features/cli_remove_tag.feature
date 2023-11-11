Feature: Remove Tags from Bookmark with CLI

  The Armaria CLI can be used to remove tags from an existing bookmark.
  
  @cli @remove_tags
  Scenario: Can remove tags
    Given the DB already has the following entries:
      | id   | parent_id | is_folder | name           | url            | description | tags              |
      | {id} | NULL      | false     | https://jho.pe | https://jho.pe | NULL        | blog, programming |
    When I run it with the following args:
      """
      remove tag [id] --tag blog --tag programming
      """
    Then the following bookmarks/folders exist:
      | id   | parent_id | is_folder | name           | url            | description | tags |
      | [id] | NULL      | false     | https://jho.pe | https://jho.pe | NULL        |      |
    And the folllowing tags exist:
      | tag |

  @cli @remove_tags
  Scenario: Bookmark must exist
    When I run it with the following args:
      """
      remove tag test --tag blog
      """
    Then the following error is returned:
      """
      Bookmark not found
      """

  @cli @remove_tags
  Scenario: Tag must exist
    Given the DB already has the following entries:
      | id   | parent_id | is_folder | name           | url            | description | tags |
      | {id} | NULL      | false     | https://jho.pe | https://jho.pe | NULL        |      |
    When I run it with the following args:
      """
      remove tag [id] --tag blog
      """
    Then the following error is returned:
      """
      Tag not found
      """
