Feature: List Parent Names

  The Armaria CLI can be used to list the parent names.

  @cli @list_parent_names
  Scenario: Can get parent names
    Given the DB already has the following entries:
      | id            | parent_id     | is_folder | name           | url            | description | tags |
      | {parent_1_id} | NULL          | true      | blogs          | NULL           | NULL        |      |
      | {parent_2_id} | [parent_1_id] | true      | programming    | NULL           | NULL        |      |
      | {id}          | [parent_2_id] | false     | https://jho.pe | https://jho.pe | NULL        |      |
    When I run it with the following args:
      """
      list parent-names [id]
      """
    Then the folllowing names are returned:
      | blogs          |
      | programming    |
      | https://jho.pe |

  @cli @list_parent_names
  Scenario: Bookmark or folder must exist
    When I run it with the following args:
      """
      list parent-names test
      """
    Then the following error is returned:
      """
      Bookmark or folder not found
      """
