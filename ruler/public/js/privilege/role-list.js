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

var jDelModal, jDelModalBody, jBtnDel, jForm, jDelId;

function initialize() {
	jDelModal = $('#modal_del_role');
	jDelModalBody = $('div.modal-body', jDelModal);
	jBtnDel = $('#btn_del_role');
	jForm = $('#form_del_role');
	jDelId = $('#txt_role_id');

	attachEvents();
}
function attachEvents() {
	jBtnDel.click(function() {
		if (jBtnDel.hasClass('disabled') || jBtnDel.attr('disabled')) {
			return false;
		}
		jBtnDel.addClass('disabled').attr('disabled', true);
        jForm.submit();
        jDelModal.modal('hide');
		return false;
	});
	jForm.submit(function() {
		jForm.ajaxSubmit({
			dataType: 'json',
			success: function(data) {
				alert(data.message);
				jBtnDel.removeClass('disabled').attr('disabled', false);
				if (data.code === 1) {
					location.reload();
				}
			}
		});
		return false;
	});
}
function preDelete(id, name) {
	jDelId.val(id);
    jDelModalBody.html([
      '<p>你确定要删除<strong>“', name, '”</strong>这个角色吗？</p>',
      '<p><strong>警告！</strong>删除角色后，所有该角色下的用户不能访问相关的资源！</p>'
      ].join(''));
    jDelModal.modal('show');
    return false;
}
window.deleteRole = preDelete;

$(function() {
	initialize();
});
})(jQuery);