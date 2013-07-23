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

  function setDisable(btn, id, disabled) {
    var url = (disabled === true ? '/admins/disable_admin/' : 
        '/admins/enable_admin/') + id;
    $.post(url, function(data) {
      var jBtn = $(btn), jLblEnabled = $('#lbl_enabled_' + id);
      if (data.code === 1) {
      	if (disabled) {
      		$('#tr_' + id).addClass('muted');
      		jBtn.attr('onclick', 
      			'return setAdminDisable(this,' + id + ',false);')
      		.html('<i class="icon-ok-circle"></i> 重新启用');
      		jLblEnabled.removeClass('badge-info').text('不可用');
      	} else {
      		$('#tr_' + id).removeClass('muted');
      		jBtn.attr('onclick', 
      			'return setAdminDisable(this,' + id + ',true);')
      		.html('<i class="icon-ban-circle"></i> 禁用');
      		jLblEnabled.addClass('badge-info').text('可用');
      	}
      }
      alert(data.message);
    }, 'json');
    return false;
  }

 window.setAdminDisable = setDisable;

})(jQuery);