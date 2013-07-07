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

var jForm, jAlert, jResName, jChkMenu, jResDescValid, jResDesc,
    jParent, jParentH, jSubmit, _timer,
    DISABLED = 'disabled';

function initElements() {
	jForm = $('#form_edit_res');
	jAlert = $('#alert');
	jResName = $('#txt_res_name');
	jChkMenu = $('#chk_is_menu');
	jParent = $('#cmb_res_parent');
	jParentH = $('#help_res_parent');
	jResDescValid = $('#txt_desc_valid');
	jResDesc = $('#txt_res_desc');
	jSubmit = $('#btn_save_res');
}

function attachEvents() {
	jForm.submit(submitHandler);
	jParent.change(function() {
		jParentH.text(jParent.find('option:selected').data('url'));
	});
	jChkMenu.click(function() {
		if (this.checked) {
			$('#txt_is_menu').val('true');
		} else {
			$('#txt_is_menu').val('false');
		}
	});
	jSubmit.click(function() {
        var resDesc = $.trim(jResDesc.val());
        jResDescValid.val(resDesc.length > 0 ? 'true' : 'false');
	});
}

function submitHandler() {
	jSubmit.button('saving').attr(DISABLED, true);
	$(this).ajaxSubmit({
		dataType: 'json',
		success: function(data) {
			if (data.code === 1) {
				showAlert('success', data.message);
				_timer = setTimeout(function() {
					hideAlert();
					clearTimeout(_timer);
					_timer = null;
				}, 6000);
			} else {
				showAlert('error', data.message);
				jSubmit.button('reset').attr(DISABLED, false);
			}
		}
	});
	return false;
}

function showAlert(type, text) {
	if (text) {
		jAlert.text(text);
	}
	if (!type) {
		type = '';
	}
    switch(type) {
    	case 'success':
    	case 'info':
    	case 'error':
    	    jAlert[0].className = 'alert alert-%s'.replace(/%s/i, type);
    	default:
    	    jAlert[0].className = 'alert';
    }
    jAlert.show();
}

function hideAlert() {
	jAlert.hide();
}

$(function() {
	initElements();
	attachEvents();
});

})(jQuery);