Feature: Add Tags to Bookmark with CLI

  The Armaria CLI can be used to add tags to an existing bookmark.

  @cli @add_tags
  Scenario: Can add tags
    Given the DB already has the following entries:
      | id   | parent_id | is_folder | name           | url            | description | tags |
      | {id} | NULL      | false     | https://jho.pe | https://jho.pe | NULL        |      |
    When I run it with the following args:
      """
      add tag [id] --tag blog --tag programming
      """
    Then the following bookmarks/folders exist:
      | id   | parent_id | is_folder | name           | url            | description | tags              |
      | [id] | NULL      | false     | https://jho.pe | https://jho.pe | NULL        | blog, programming |
    And the folllowing tags exist:
      | tag         |
      | blog        |
      | programming |

  @cli @add_tags
  Scenario: Bookmark must exist
    When I run it with the following args:
      """
      add tag test --tag blog
      """
    Then the following error is returned:
      """
      Bookmark not found
      """

  @cli @add_tags
  Scenario: Tags must be in the char range [A-Z][a-z][0-9]-_
    Given the DB already has the following entries:
      | id   | parent_id | is_folder | name           | url            | description | tags |
      | {id} | NULL      | false     | https://jho.pe | https://jho.pe | NULL        |      |
    When I run it with the following args:
      """
      add tag [id] --tag ?
      """
    Then the following error is returned:
      """
      Tag has invalid chars
      """

  @cli @add_tags
  Scenario: Tag must be at most 128 chars
    Given the DB already has the following entries:
      | id   | parent_id | is_folder | name           | url            | description | tags |
      | {id} | NULL      | false     | https://jho.pe | https://jho.pe | NULL        |      |
    When I run it with the following args:
      """
      add tag [id] --tag %repeat:x:129%
      """
    Then the following error is returned:
      """
      Tag too long
      """

  @cli @add_tags
  Scenario: Can have at most 24 tags
    Given the DB already has the following entries:
      | id   | parent_id | is_folder | name           | url            | description | tags |
      | {id} | NULL      | false     | https://jho.pe | https://jho.pe | NULL        | blog |
    When I run it with the following args:
      """
      add tag [id] %repeat: --tag "[uuid]":24%
      """
    Then the following error is returned:
      """
      Too many tags applied to bookmark
      """

  @cli @add_tags
  Scenario: Tags must be unique
    Given the DB already has the following entries:
      | id   | parent_id | is_folder | name           | url            | description | tags |
      | {id} | NULL      | false     | https://jho.pe | https://jho.pe | NULL        |      |
    When I run it with the following args:
      """
      add tag [id] --tag blog --tag blog
      """
    Then the following error is returned:
      """
      Tags must be unique
      """

  @cli @add_tags
  Scenario: Cannot add same tag again
    Given the DB already has the following entries:
      | id   | parent_id | is_folder | name           | url            | description | tags |
      | {id} | NULL      | false     | https://jho.pe | https://jho.pe | NULL        | blog |
    When I run it with the following args:
      """
      add tag [id] --tag blog
      """
    Then the following error is returned:
      """
      Tags must be unique
      """
