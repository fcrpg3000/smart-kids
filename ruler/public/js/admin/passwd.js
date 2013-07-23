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

/**
 * @author king4go fcrpg3000 (fcrpg2005 At gmail.com)
 * @since 1.0
 */
(function($) {

var jForm, jOldpwd, jNewpwd1, jNewpwd2, jOldpwdC, jNewpwd1C, jNewpwd2C,
    jOldpwdHelper, jNewpwd1Helper, jNewpwd2Helper, jAlert, jPasswd, 
    CG_SEL = 'div.control-group';

function initElements() {
    jForm = $('#form_chpwd');
    jOldpwd = $('#pwd_oldpwd');
    jNewpwd1 = $('#pwd_newpwd1');
    jNewpwd2 = $('#pwd_newpwd2');
    jOldpwdHelper = $('#help_oldpwd');
    jNewpwd1Helper = $('#help_newpwd1');
    jNewpwd2Helper = $('#help_newpwd2');

    jOldpwdC = jOldpwd.parents(CG_SEL);
    jNewpwd1C = jNewpwd1.parents(CG_SEL);
    jNewpwd2C = jNewpwd2.parents(CG_SEL);
    jAlert = $('#message_tip');
    jPasswd = $('#btn_passwd');
}

function attachEvents() {
	jOldpwd.focus(function() {
		if (jOldpwdC.hasClass('error')) {
			infoCtrlGroup(jOldpwdC);
			jOldpwdHelper.hide();
		}
		hideAlert();
	});
	jNewpwd1.focus(function() {
        if (jNewpwd1C.hasClass('error')) {
            infoCtrlGroup(jNewpwd1C);
            jNewpwd1Helper.hide();
        }
        hideAlert();
    });
    jNewpwd2.focus(function() {
        if (jNewpwd2C.hasClass('error')) {
            infoCtrlGroup(jNewpwd2C);
            jNewpwd2Helper.hide();
        }
        hideAlert();
    });
    jPasswd.click(function() {
        var oldPwd = $.trim(jOldpwd.val()),
                newPwd1 = $.trim(jNewpwd1.val()),
                newPwd2 = $.trim(jNewpwd2.val());
        if (!oldPwd.length) {
            errorCtrlGroup(jOldpwdC);
            jOldpwdHelper.show();
            jOldPwd.focus();
            return false;
        }
        if (!newPwd1.length) {
            errorCtrlGroup(jNewpwd1C);
            jNewpwd1Helper.show();
            jNewpwd1.focus();
            return false;
        }
        if (!newPwd2.length || newPwd1 !== newPwd2) {
            errorCtrlGroup(jNewpwd2C);
            jNewpwd2Helper.show();
            return false;
        }
        return true;
    });
}

function errorCtrlGroup(jcg) {
    jcg.addClass('error');
}

function infoCtrlGroup(jcg) {
	jcg.removeClass('error');
}

function hideAlert() {
	if (jAlert.hasClass('alert-error'))
        jAlert.removeClass('alert-error').addClass('alert-info').hide();
}

$(function() {
	initElements();
	attachEvents();
	jForm.ajaxForm({
		dataType: 'json',
		success: function(data) {
			if (data.code === 1) {
                alert(data.message);
                location.href = "/";
                return;
            } else {
                jAlert.addClass('alert-error').removeClass('alert-info').text(data.message).show();
            }
		}
	});
});

})(jQuery);