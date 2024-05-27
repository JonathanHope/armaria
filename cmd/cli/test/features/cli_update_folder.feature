Feature: Update Folders with CLI

  The Armaria CLI can be used to update an existing folder.

  @cli @update_folder
  Scenario: Can update folder name
    Given the DB already has the following entries:
      | id   | parent_id | is_folder | name  | url  | description | tags |
      | {id} | NULL      | true      | blogs | NULL | NULL        |      |
    When I run it with the following args:
      """
      update folder [id] --name new
      """
    Then the following bookmarks/folders exist:
      | id   | parent_id | is_folder | name | url  | description | tags |
      | [id] | NULL      | true      | new  | NULL | NULL        |      |

  @cli @update_folder
  Scenario: Can move folder
    Given the DB already has the following entries:
      | id          | parent_id | is_folder | name  | url  | description | tags |
      | {parent_id} | NULL      | true      | tech  | NULL | NULL        |      |
      | {id}        | NULL      | true      | blogs | NULL | NULL        |      |
    When I run it with the following args:
      """
      update folder [id] --folder [parent_id]
      """
    Then the following bookmarks/folders exist:
      | id          | parent_id   | is_folder | name  | url  | description | tags |
      | [parent_id] | NULL        | true      | tech  | NULL | NULL        |      |
      | [id]        | [parent_id] | true      | blogs | NULL | NULL        |      |

  @cli @update_folder
  Scenario: Can remove parent folder
    Given the DB already has the following entries:
      | id          | parent_id   | is_folder | name  | url  | description | tags |
      | {parent_id} | NULL        | true      | tech  | NULL | NULL        |      |
      | {id}        | [parent_id] | true      | blogs | NULL | NULL        |      |
    When I run it with the following args:
      """
      update folder [id] --no-folder
      """
    Then the following bookmarks/folders exist:
      | id          | parent_id | is_folder | name  | url  | description | tags |
      | [parent_id] | NULL      | true      | tech  | NULL | NULL        |      |
      | [id]        | NULL      | true      | blogs | NULL | NULL        |      |

  @cli @update_folder
  Scenario: Can move folder to start of list
    Given the DB already has the following entries:
      | id     | parent_id | is_folder | name  | url  | description | tags |
      | {id_1} | NULL      | true      | one   | NULL | NULL        |      |
      | {id_2} | NULL      | true      | two   | NULL | NULL        |      |
      | {id_3} | NULL      | true      | three | NULL | NULL        |      |
    When I run it with the following args:
      """
      update folder [id_2] --before [id_1] 
      """
    Then the following bookmarks/folders exist:
      | id     | parent_id | is_folder | name  | url  | description | tags |
      | [id_2] | NULL      | true      | two   | NULL | NULL        |      |
      | [id_1] | NULL      | true      | one   | NULL | NULL        |      |
      | [id_3] | NULL      | true      | three | NULL | NULL        |      |

  @cli @update_folder
  Scenario: Can move folder to end of list
    Given the DB already has the following entries:
      | id     | parent_id | is_folder | name  | url  | description | tags |
      | {id_1} | NULL      | true      | one   | NULL | NULL        |      |
      | {id_2} | NULL      | true      | two   | NULL | NULL        |      |
      | {id_3} | NULL      | true      | three | NULL | NULL        |      |
    When I run it with the following args:
      """
      update folder [id_2] --after [id_3] 
      """
    Then the following bookmarks/folders exist:
      | id     | parent_id | is_folder | name  | url  | description | tags |
      | [id_1] | NULL      | true      | one   | NULL | NULL        |      |
      | [id_3] | NULL      | true      | three | NULL | NULL        |      |
      | [id_2] | NULL      | true      | two   | NULL | NULL        |      |

  @cli @update_folder
  Scenario: Can move folder to midddle of list
    Given the DB already has the following entries:
      | id     | parent_id | is_folder | name  | url  | description | tags |
      | {id_1} | NULL      | true      | one   | NULL | NULL        |      |
      | {id_2} | NULL      | true      | two   | NULL | NULL        |      |
      | {id_3} | NULL      | true      | three | NULL | NULL        |      |
    When I run it with the following args:
      """
      update folder [id_1] --after [id_2] --before [id_3] 
      """
    Then the following bookmarks/folders exist:
      | id     | parent_id | is_folder | name  | url  | description | tags |
      | [id_2] | NULL      | true      | two   | NULL | NULL        |      |
      | [id_1] | NULL      | true      | one   | NULL | NULL        |      |
      | [id_3] | NULL      | true      | three | NULL | NULL        |      |

  @cli @update_folder
  Scenario: Name must be at most 2048 chars
    Given the DB already has the following entries:
   | id   | parent_id | is_folder | name  | url  | description | tags |
   | {id} | NULL      | true      | blogs | NULL | NULL        |      |
  When I run it with the following args:
    """
    update folder [id] --name %repeat:x:2049%
    """
  Then the following error is returned:
    """
    Name too long
    """

  @cli @update_folder
  Scenario: Name must be at least 1 char
    Given the DB already has the following entries:
      | id   | parent_id | is_folder | name  | url  | description | tags |
      | {id} | NULL      | true      | blogs | NULL | NULL        |      |
    When I run it with the following args:
      """
      update folder [id] --name ""
      """
    Then the following error is returned:
      """
      Name too short
      """

  @cli @update_folder
  Scenario: Name must be at most 2048 chars
    Given the DB already has the following entries:
      | id   | parent_id | is_folder | name  | url  | description | tags |
      | {id} | NULL      | true      | blogs | NULL | NULL        |      |
    When I run it with the following args:
      """
      update folder [id] --name %repeat:x:2049%
      """
    Then the following error is returned:
      """
      Name too long
      """
      
  @cli @update_folder
  Scenario: Parent folder must exist
    Given the DB already has the following entries:
      | itd  | parent_id | is_folder | name  | url  | description | tags |
      | {id} | NULL      | true      | blogs | NULL | NULL        |      |
    When I run it with the following args:
      """
      update folder [id] --folder test
      """
    Then the following error is returned:
      """
      Folder not found
      """

  @cli @update_folder
  Scenario: Cannot move and remove folder at the same time
    Given the DB already has the following entries:
      | id          | parent_id | is_folder | name  | url  | description | tags |
      | {parent_id} | NULL      | true      | tech  | NULL | NULL        |      |
      | {id}        | NULL      | true      | blogs | NULL | NULL        |      |
  When I run it with the following args:
    """
    update folder [id] --folder [parent_id] --no-folder
    """
  Then the following error is returned:
    """
    Arguments folder and no-folder are mutually exclusive
    """
    
  @cli @update_folder
  Scenario: Folder must exist
    When I run it with the following args:
      """
      update folder test --name "new"
      """
    Then the following error is returned:
      """
      Folder not found
      """

  @cli @update_folder
  Scenario: At least one update is required
    Given the DB already has the following entries:
      | id   | parent_id | is_folder | name  | url  | description | tags |
      | {id} | NULL      | true      | blogs | NULL | NULL        |      |
    When I run it with the following args:
      """
      update folder [id]
      """
    Then the following error is returned:
      """
      At least one update is required
      """
