/* 
 * Copyright (C) 2012-2013 king4go authors All rights reserved.
 *
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements. See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 *
 *           http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

(function($) {
  var jAddRole, jRemoveRole, jAllRoles, jRoles,
      DISABLED = 'disabled', SELECTED = 'selected',
      SELECTED_OPTION = 'option:selected';

  /* init all html elements and variables when doc init. */
  function initElements() {
    jAddRole = $('#btn_add_role');
    jRemoveRole = $('#btn_remove_role');
    jAllRoles = $('#cmb_all_roles');
    jRoles = $('#cmb_roles');

    $(SELECTED_OPTION, jAllRoles).attr(SELECTED, false);
    $(SELECTED_OPTION, jRoles).attr(SELECTED, false);
  }
  /* attach all events binding when doc init.*/
  function attachEvents() {
    jAllRoles.on('change click', allRoleChangeHandler);
    jRoles.on('change click', rolesChangeHandler);
    jAddRole.click(addRolesHandler);
    jRemoveRole.click(removeRolesHandler);
  }

  /* Add roles events handler. */
  function addRolesHandler() {
    var $options = $(SELECTED_OPTION, jAllRoles);
      if (!$options.length) {
        return false;
      }
      var $removed = $options.remove();
      $removed.appendTo(jRoles).attr(SELECTED, false);
      if (!$(SELECTED_OPTION, jAllRoles).length) {
        $(this).addClass(DISABLED).attr(DISABLED, true);
      }
      return false;
  }
  /* Remove roles events handler. */
  function removeRolesHandler() {
    var $options = $(SELECTED_OPTION, jRoles);
      if (!$options.length) {
        return false;
      }
      var $removed = $options.remove();
      $removed.appendTo(jAllRoles).attr(SELECTED, false);
      if (!$(SELECTED_OPTION, jRoles).length) {
        $(this).addClass(DISABLED).attr(DISABLED, true);
      }
      return false;
  }

  function allRoleChangeHandler() {
    if ($(this).val() != null) {
      jAddRole.attr(DISABLED, false).removeClass(DISABLED);
    } else {
      jAddRole.attr(DISABLED, true).addClass(DISABLED);
    }
    return false;
  }
  function rolesChangeHandler() {
    if ($(this).val() != null) {
      jRemoveRole.attr(DISABLED, false).removeClass(DISABLED);
    } else {
      jRemoveRole.attr(DISABLED, true).addClass(DISABLED);
    }
    return false;
  }
  $(function() {
    initElements();
    attachEvents();
  });
})(jQuery);