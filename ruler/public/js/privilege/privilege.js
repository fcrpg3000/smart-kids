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

$(function() {
	$('input[name=mainResource]').click(function() {
		$('input[name=' + $(this).data('subName') + ']')
		  .attr('checked', this.checked ? true : false);
	});
	$('input[name^=subResource]').click(function() {
		var $this = $(this), name = $this.attr('name'), jParent = $($this.data('parent'));
		if (this.checked) {
			jParent[0].checked = true;
		} else {
			if (!$('input[name=' + name + ']:checked').length) {
                jParent[0].checked = false;
			}
		}
	});
	$('#ch_res_all').click(function() {
		$('input:checkbox').attr('checked', this.checked ? true: false);
	});
});

})(jQuery);