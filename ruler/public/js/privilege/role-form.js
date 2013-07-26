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
    var ERR_CLS = 'error', DISABLED = 'disabled', CG_SELECTOR = 'div.control-group',
            HELP_SELECTOR = 'span.help-inline', SOURCE_CONTENT = 'srcContent', _timer,
            jForm, jName, jCode, jUri, jMessageTip,
            jNameCG, jCodeCG, jUriCG,
            jNameHelper, jCodeHelper, jUriHelper,
            jChkMenu, jTxtMenu, jSubmit, jReset;
    function _initElements() {
        jForm = $('#form_edit_role');
        jName = $('#txt_role_name');
        jCode = $('#txt_role_code');
        jMessageTip = $('#message_tip');

        jNameCG = jName.parents(CG_SELECTOR);
        jCodeCG = jCode.parents(CG_SELECTOR);

        jNameHelper = jName.next(HELP_SELECTOR);
        jCodeHelper = jCode.next(HELP_SELECTOR);

        jChkMenu = $('#chk_is_enabled');
        jTxtMenu = $('#txt_is_enabled');
        jSubmit = $('#btn_save');
        jReset = $('#btn_reset');
    }

    function _ctrlFocus(obj) {
        if (obj.cg.hasClass(ERR_CLS)) {
            obj.cg.removeClass(ERR_CLS);
            var content = obj.helper.data(SOURCE_CONTENT);
            if (content) {
                obj.helper.show().text(content);
            }
        }
    }

    function _initEvents() {
        jChkMenu.click(function() {
            if (this.checked) {
                jTxtMenu.val('true');
            } else {
                jTxtMenu.val('false');
            }
        });

        jName.focus(function() {
            _ctrlFocus({
                cg: jNameCG,
                helper: jNameHelper
            });
        });

        jSubmit.click(function() {
            var name = $.trim(jName.val()),
                    code = $.trim(jCode.val());
            _disableForm();

            if (!name.length) {
                jNameCG.addClass(ERR_CLS);
                jNameHelper.show();
                _enableForm();
                return false;
            }
            if (!code.length) {
                jCodeCG.addClass(ERR_CLS);
                jCodeHelper.show();
                _enableForm();
                return false;
            }
           
            return true;
        });
    }
    // 解绑元素事件
    function _detachEvents() {
        jSubmit.off();
        jReset.off();
        jChkMenu.off();
        jForm.off();
    }
    // 禁用表单元素
    function _disableForm() {
        jReset.attr(DISABLED, true).addClass(DISABLED);

    }
    // 
    function _enableForm() {
        jSubmit.attr(DISABLED, false).removeClass(DISABLED);
        jReset.attr(DISABLED, false).removeClass(DISABLED);
    }
    function _initForm() {

        jForm.ajaxForm({
            dataType: 'json',
            success: function(data) {
                var isSuccess = data.code === 1;
                jMessageTip.show().text(data.message).addClass(isSuccess ?
                        'text-success' : 'text-error');
                _timer = setTimeout(function() {
                    jMessageTip.removeClass('text-success')
                            .removeClass('text-error').hide();

                    clearTimeout(_timer);
                    _timer = null;

                    if (isSuccess) {
                        location.href = '/privilege/role_list';
                    } else {
                        _enableForm();
                    }
                }, 3000);
            }
        });
    }
    function _initInternal() {
        _initElements();
        _initEvents();
        _initForm();
    }

    $(function() {
        _initInternal();

        $(document).on('unload', function() {
            _detachEvents();
        });
    });

})(jQuery);