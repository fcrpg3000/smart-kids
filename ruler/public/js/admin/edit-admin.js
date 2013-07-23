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
  var jAddRole, jRemoveRole, jAllRoles, jRoles, jForm, jAid,
      jAdminName, jEmpName, jEmpNo, jAdminNameC, jEmpNameC, jEmpNoC, 
      jAlert, jSubmit, jModal, isModify = 0, checkedName = {}, errorName = {},
      DISABLED = 'disabled', SELECTED = 'selected', CTRL_GROUP = 'div.control-group',
      SELECTED_OPTION = 'option:selected';

  /* init all html elements and variables when doc init. */
  function initElements() {
    jForm = $('#form_edit_admin');
    jAlert = $('#message_tip');
    jAid = $('#txt_admin_id');
    jAdminName = $('#txt_admin_name');
    jEmpName = $('#txt_emp_name');
    jEmpNo = $('#txt_emp_no');
    jAdminNameC = jAdminName.parents(CTRL_GROUP);
    jEmpNameC = jEmpName.parents(CTRL_GROUP);
    jEmpNoC = jEmpNo.parents(CTRL_GROUP);
    jAddRole = $('#btn_add_role');
    jRemoveRole = $('#btn_remove_role');
    jAllRoles = $('#cmb_all_roles');
    jRoles = $('#cmb_roles');
    jSubmit = $('#btn_save_admin').addClass(DISABLED).attr(DISABLED, true);
    jModal = $('#confirm_modal');

    if (jAid.length && $.trim(jAid.val()) != "") {
      isModify = 1;
      oldName = jAdminName.data('oldName');
    }

    $(SELECTED_OPTION, jAllRoles).attr(SELECTED, false);
    $(SELECTED_OPTION, jRoles).attr(SELECTED, false);
  }
  /* attach all events binding when doc init.*/
  function attachEvents() {
    jAllRoles.on('change click', allRoleChangeHandler);
    jRoles.on('change click', rolesChangeHandler);
    jAddRole.click(addRolesHandler);
    jRemoveRole.click(removeRolesHandler);
    jForm.submit(formSubmitHandler);

    /* jAdminName focus and blur events binding. */
    jAdminName.focus(function() {
      hideAlert();
      normalCtrlGroup(jAdminNameC);
      return false;
    }).blur(function() {
      var $this = $(this), val = $.trim(this.value);
      if (!val.length) {
        showAlert('请输入管理员用户名称！');
        errorCtrlGroup(jAdminNameC);
        return false;
      }
      if ((isModify && oldName != val) || !isModify) {
        if (checkedName[val]) {
          successCtrlGroup(jAdminNameC);
          enableSubmit();
          return false;
        }
        if (errorName[val]) {
          showAlert();
          errorCtrlGroup(jAdminNameC);
          disableSubmit();
          return false;
        }
        // admin name changed, need verify
        checkName(val);
      } else {
        successCtrlGroup(jAdminNameC);
        enableSubmit();
      }
      return false;
    });
    jSubmit.click(function() {
      jSubmit.button('saving');
      var admName = $.trim(jAdminName.val()),
          empName = $.trim(jEmpName.val()),
          empNo = $.trim(jEmpNo.val());
      if (!admName.length) {
        showAlert('请输入管理员用户名称！');
        errorCtrlGroup(jAdminNameC);
        jSubmit.button('reset');
        return false;
      }
      jForm.submit();
      return false;
    });
    $('#btn_goon').click(function() {
      jForm.resetForm();
      jModal.modal('hide');
      location.reload();
      return false;
    });
  }

  function checkName(val, fn1, fn2) {
    $.getJSON('/admins/check_admin_name/' + val, 
      function(data) {
        if (data.code === 1) {
          hideAlert();
          successCtrlGroup(jAdminNameC);
          checkedName[val] = 1;
          enableSubmit();
        } else {
          showAlert(data.message);
          errorCtrlGroup(jAdminNameC);
          errorName[val] = 1;
          disableSubmit();
        }
      });
  }

  function errorCtrlGroup(jCtrlGroup) {
    jCtrlGroup.removeClass('success').addClass('error');
  }

  function normalCtrlGroup(jCtrlGroup) {
    if (jCtrlGroup.hasClass('success')) {
      return;
    }
    jCtrlGroup.removeClass('error');
  }

  function successCtrlGroup(jCtrlGroup) {
    jCtrlGroup.removeClass('error').addClass('success');
  }

  function enableSubmit() {
    jSubmit.removeClass(DISABLED).attr(DISABLED, false);
  }
  function disableSubmit() {
    jSubmit.addClass(DISABLED).attr(DISABLED, true);
  }

  function showAlert(text) {
    if (text) {
      jAlert.html(text);
    }
    jAlert.show();
  }

  function hideAlert() {
    jAlert.hide();
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
  function formSubmitHandler() {
    $(this).ajaxSubmit({
      dataType: 'json',
      success: function(data) {
        if (data.code !== 1) {
          jSubmit.button('reset');
          var tips = ['<strong>', data.message, '</strong>'];
          if (data.values) {
            for (var k in data.values) {
              tips.push(['<p>', data.values[k], '</p>'].join(''));
              errorCtrlGroup($('input[name="' + k + '"]').parents(CTRL_GROUP))
            }
          }
          showAlert(tips.join(''));
        } else {
          $('div.modal-body', jModal).html('<p>' + data.message + '</p>');
          jModal.modal('show');
        }
      }
    });
    return false;
  }
  
  $(function() {
    initElements();
    attachEvents();
  });
})(jQuery);