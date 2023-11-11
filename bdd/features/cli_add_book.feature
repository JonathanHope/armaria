Feature: Add Bookmarks with CLI

  The Armaria CLI can be used to add a new bookmark.

  @cli @add_book
  Scenario: Can add a bookmark
    When I run it with the following args:
      """
      add book https://jho.pe
      """
    Then the following bookmarks/folders exist: 
      | id   | parent_id | is_folder | name           | url            | description | tags |
      | {id} | NULL      | false     | https://jho.pe | https://jho.pe | NULL        |      |

  @cli @add_book
  Scenario: Can add a boookmark to a folder
    Given the DB already has the following entries:
      | id          | parent_id | is_folder | name  | url  | description | tags |
      | {parent_id} | NULL      | true      | blogs | NULL | NULL        |      |
    When I run it with the following args:
      """
      add book https://jho.pe --folder [parent_id]
      """
    Then the following bookmarks/folders exist:
      | id          | parent_id   | is_folder | name           | url            | description | tags |
      | [parent_id] | NULL        | true      | blogs          | NULL           | NULL        |      |
      | {id}        | [parent_id] | false     | https://jho.pe | https://jho.pe | NULL        |      |

  @cli @add_book
  Scenario: Can add a bookmark with a name
    When I run it with the following args:
      """
      add book https://jho.pe --name "The Flat Field"
      """
    Then the following bookmarks/folders exist:
      | id   | parent_id | is_folder | name           | url            | description | tags |
      | {id} | NULL      | false     | The Flat Field | https://jho.pe | NULL        |      |

  @cli @add_book
  Scenario: Can add a bookmark with a description
    When I run it with the following args:
      """
      add book https://jho.pe --description "The blog of Jonathan Hope."
      """
    Then the following bookmarks/folders exist:
      | id   | parent_id | is_folder | name           | url            | description                | tags |
      | {id} | NULL      | false     | https://jho.pe | https://jho.pe | The blog of Jonathan Hope. |      |

  @cli @add_book
  Scenario: Can add a bookmark with tags
    When I run it with the following args:
      """
      add book https://jho.pe --tag "blog" --tag "programming"
      """
    Then the following bookmarks/folders exist:
      | id   | parent_id | is_folder | name           | url            | description | tags              |
      | {id} | NULL      | false     | https://jho.pe | https://jho.pe | NULL        | blog, programming |
    And the folllowing tags exist:
      | tag         |
      | blog        |
      | programming |

  @cli @add_book
  Scenario: Folder must exist
    When I run it with the following args:
      """
      add book https://jho.pe --folder test
      """
    Then the following error is returned:
      """
      Folder not found
      """

  @cli @add_book
  Scenario: URL must be at least 1 char
    When I run it with the following args:
      """
      add book ""
      """
    Then the following error is returned:
      """
      URL too short
      """

  @cli @add_book
  Scenario: URL must be at most 2048 chars
    When I run it with the following args:
      """
      add book %repeat:x:2049%
      """
    Then the following error is returned:
      """
      URL too long
      """

  @cli @add_book
  Scenario: Name must be at most 2048 chars
    When I run it with the following args:
      """
      add book https://jho.pe --name %repeat:x:2049%
      """
    Then the following error is returned:
      """
      Name too long
      """

  @cli @add_book
  Scenario: Description must be at leat 1 char
    When I run it with the following args:
      """
      add book https://jho.pe --description ""
      """
    Then the following error is returned:
      """
      Description too short
      """

  @cli @add_book
  Scenario: Description must be at most 4096 chars
    When I run it with the following args:
      """
      add book https://jho.pe --description %repeat:x:4097%
      """
    Then the following error is returned:
      """
      Description too long
      """

  @cli @add_book
  Scenario: Tag must be at most 128 chars
    When I run it with the following args:
      """
      add book https://jho.pe --tag %repeat:x:129%
      """
    Then the following error is returned:
      """
      Tag too long
      """

  @cli @add_book
  Scenario: Tags must be unique
    When I run it with the following args:
      """
      add book https://jho.pe --tag blog --tag blog
      """
    Then the following error is returned:
      """
      Tags must be unique
      """

  @cli @add_book
  Scenario: Can have at most 24 tags
    When I run it with the following args:
      """
      add book https://jho.pe %repeat: --tag "[uuid]":25%
      """
    Then the following error is returned:
      """
      Too many tags applied to bookmark
      """
      
  @cli @add_book
  Scenario: Tags must be in the char range [A-Z][a-z][0-9]-_
    When I run it with the following args:
      """
      add book https://jho.pe --tag ?
      """
    Then the following error is returned:
      """
      Tag has invalid chars
      """
