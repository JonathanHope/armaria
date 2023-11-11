Feature: Update Bookmarks with CLI

  The Armaria CLI can be used to update an existing bookmark.

  @cli @update_book
  Scenario: Can update bookmark URL
    Given the DB already has the following entries:
      | id   | parent_id | is_folder | name           | url            | description | tags |
      | {id} | NULL      | false     | https://jho.pe | https://jho.pe | NULL        |      |
    When I run it with the following args:
      """
      update book [id] --url https://theflatfield.net
      """
    Then the following bookmarks/folders exist:
      | id   | parent_id | is_folder | name           | url                      | description | tags |
      | [id] | NULL      | false     | https://jho.pe | https://theflatfield.net | NULL        |      |

  @cli @update_book
  Scenario: Can update bookmark name
    Given the DB already has the following entries:
      | id   | parent_id | is_folder | name           | url            | description | tags |
      | {id} | NULL      | false     | https://jho.pe | https://jho.pe | NULL        |      |
    When I run it with the following args:
      """
      update book [id] --name "The Flat Field"
      """
    Then the following bookmarks/folders exist:
      | id   | parent_id | is_folder | name           | url            | description | tags |
      | [id] | NULL      | false     | The Flat Field | https://jho.pe | NULL        |      |

  @cli @update_book
  Scenario: Can update bookmark description
    Given the DB already has the following entries:
      | id   | parent_id | is_folder | name           | url            | description | tags |
      | {id} | NULL      | false     | https://jho.pe | https://jho.pe | NULL        |      |
    When I run it with the following args:
      """
      update book [id] --description "The blog of Jonathan Hope."
      """
    Then the following bookmarks/folders exist:
      | id   | parent_id | is_folder | name           | url            | description                | tags |
      | [id] | NULL      | false     | https://jho.pe | https://jho.pe | The blog of Jonathan Hope. |      |

  @cli @update_book
  Scenario: Can remove bookmark description
    Given the DB already has the following entries:
      | id   | parent_id | is_folder | name           | url            | description                | tags |
      | {id} | NULL      | false     | https://jho.pe | https://jho.pe | The blog of Jonathan Hope. |      |
    When I run it with the following args:
      """
      update book [id] --no-description
      """
    Then the following bookmarks/folders exist:
      | id   | parent_id | is_folder | name           | url            | description | tags |
      | [id] | NULL      | false     | https://jho.pe | https://jho.pe | NULL        |      |

  @cli @update_book
  Scenario: Can move bookmark
    Given the DB already has the following entries:
      | id          | parent_id | is_folder | name           | url            | description | tags |
      | {parent_id} | NULL      | true      | blogs          | NULL           | NULL        |      |
      | {id}        | NULL      | false     | https://jho.pe | https://jho.pe | NULL        |      |
    When I run it with the following args:
      """
      update book [id] --folder [parent_id]
      """
    Then the following bookmarks/folders exist:
      | id          | parent_id   | is_folder | name           | url            | description | tags |
      | [parent_id] | NULL        | true      | blogs          | NULL           | NULL        |      |
      | [id]        | [parent_id] | false     | https://jho.pe | https://jho.pe | NULL        |      |

  @cli @update_book
  Scenario: Can remove parent folder
    Given the DB already has the following entries:
      | id          | parent_id   | is_folder | name           | url            | description | tags |
      | {parent_id} | NULL        | true      | blogs          | NULL           | NULL        |      |
      | {id}        | [parent_id] | false     | https://jho.pe | https://jho.pe | NULL        |      |
    When I run it with the following args:
      """
      update book [id] --no-folder
      """
    Then the following bookmarks/folders exist:
      | id          | parent_id | is_folder | name           | url            | description | tags |
      | [parent_id] | NULL      | true      | blogs          | NULL           | NULL        |      |
      | [id]        | NULL      | false     | https://jho.pe | https://jho.pe | NULL        |      |

  @cli @update_book
  Scenario: Bookmark must exist
    When I run it with the following args:
      """
      update book test --name "The Flat Field"
      """
    Then the following error is returned:
      """
      Bookmark not found
      """

  @cli @update_book
  Scenario: At least one update is required
    Given the DB already has the following entries:
      | id   | parent_id | is_folder | name           | url            | description | tags |
      | {id} | NULL      | false     | https://jho.pe | https://jho.pe | NULL        |      |
    When I run it with the following args:
      """
      update book [id]
      """
    Then the following error is returned:
      """
      At least one update is required
      """

  @cli @update_book
  Scenario: Name must be at least 1 char
    Given the DB already has the following entries:
      | id   | parent_id | is_folder | name           | url            | description | tags |
      | {id} | NULL      | false     | https://jho.pe | https://jho.pe | NULL        |      |
    When I run it with the following args:
      """
      update book [id] --name ""
      """
    Then the following error is returned:
      """
      Name too short
      """

  @cli @update_book
  Scenario: Name must be at most 2048 chars
    Given the DB already has the following entries:
      | id   | parent_id | is_folder | name           | url            | description | tags |
      | {id} | NULL      | false     | https://jho.pe | https://jho.pe | NULL        |      |
    When I run it with the following args:
      """
      update book [id] --name %repeat:x:2049%
      """
    Then the following error is returned:
      """
      Name too long
      """

  @cli @update_book
  Scenario: URL must be at least 1 char
    Given the DB already has the following entries:
      | id   | parent_id | is_folder | name           | url            | description | tags |
      | {id} | NULL      | false     | https://jho.pe | https://jho.pe | NULL        |      |
    When I run it with the following args:
      """
      update book [id] --url ""
      """
    Then the following error is returned:
      """
      URL too short
      """
      
  @cli @update_book
  Scenario: URL must be at most 2048 chars
    Given the DB already has the following entries:
      | id   | parent_id | is_folder | name           | url            | description | tags |
      | {id} | NULL      | false     | https://jho.pe | https://jho.pe | NULL        |      |
    When I run it with the following args:
      """
      update book [id] --url %repeat:x:2049%
      """
    Then the following error is returned:
      """
      URL too long
      """

  @cli @update_book
  Scenario: Description must be at least 1 char
    Given the DB already has the following entries:
      | id   | parent_id | is_folder | name           | url            | description | tags |
      | {id} | NULL      | false     | https://jho.pe | https://jho.pe | NULL        |      |
    When I run it with the following args:
      """
      update book [id] --description ""
      """
    Then the following error is returned:
      """
      Description too short
      """

  @cli @update_book
  Scenario: Description must be at most 4096 chars
    Given the DB already has the following entries:
      | id   | parent_id | is_folder | name           | url            | description | tags |
      | {id} | NULL      | false     | https://jho.pe | https://jho.pe | NULL        |      |
    When I run it with the following args:
      """
      update book [id] --description %repeat:x:4097%
      """
    Then the following error is returned:
      """
      Description too long
      """

  @cli @update_book
  Scenario: Cannot update and remove description at the same time
    Given the DB already has the following entries:
      | id   | parent_id | is_folder | name           | url            | description | tags |
      | {id} | NULL      | false     | https://jho.pe | https://jho.pe | NULL        |      |
    When I run it with the following args:
      """
      update book [id] --description "The blog of Jonathan Hope." --no-description
      """
    Then the following error is returned:
      """
      Arguments description and no-description are mutually exclusive
      """

  @cli @update_book
  Scenario: Cannot move and remove folder at the same time
    Given the DB already has the following entries:
      | id          | parent_id | is_folder | name           | url            | description | tags |
      | {parent_id} | NULL      | true      | tech           | NULL           | NULL        |      |
      | {id}        | NULL      | false     | https://jho.pe | https://jho.pe | NULL        |      |
    When I run it with the following args:
      """
      update book [id] --folder [parent_id] --no-folder
      """
    Then the following error is returned:
      """
      Arguments folder and no-folder are mutually exclusive
      """

  @cli @update_book
  Scenario: Parent folder must exist
    Given the DB already has the following entries:
      | id   | parent_id | is_folder | name           | url            | description | tags |
      | {id} | NULL      | false     | https://jho.pe | https://jho.pe | NULL        |      |
    When I run it with the following args:
      """
      update book [id] --folder test
      """
    Then the following error is returned:
      """
      Folder not found
      """
